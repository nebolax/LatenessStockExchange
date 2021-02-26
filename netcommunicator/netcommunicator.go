package netcommunicator

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
)

//OutcomingMessage is a struct
type OutcomingMessage struct {
	Type        string  `json:"type"`
	OffersCount int     `json:"offersCount"`
	StockPrice  float64 `json:"stockPrice"`
}

//incomingMessage is struct too
type incomingMessage struct {
	OfferType string `json:"offerType"`
}
type graphPageSetup struct {
	Type      string    `json:"type"`
	History   []float64 `json:"history"`
	CurOffers int       `json:"offers"`
}

var (
	clients   = make(map[*websocket.Conn]int)
	broadcast = make(chan int)
	key       = []byte("super-secret-key")
	store     = sessions.NewCookieStore(key)
)

func getUserID(w http.ResponseWriter, r *http.Request) int {
	session, _ := store.Get(r, "user-info")
	id, ok := session.Values["userid"].(int)

	if !ok || (id == 0) {
		session.Values["userid"] = 0
		session.Save(r, w)
		return 0
	}

	return id
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func mhandler(w http.ResponseWriter, req *http.Request) {
	html, err := ioutil.ReadFile("./templates/graph-page.html")
	checkerr(err)
	w.Write(html)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	clients[ws] = id
	defer ws.Close()

	if calc := pricesmonitor.CalcByID(id); calc != nil {
		fmt.Printf("Id: %d, Offs: %d, st: %#v\n", calc.ID, calc.CurHandler.CurOffers, calc.History)
		ws.WriteJSON(graphPageSetup{"setup", calc.History, calc.CurHandler.CurOffers})
	}

	//Reading data from clients
	for {
		var msg incomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		var inc int
		if msg.OfferType == "sell" {
			inc = -1
		} else if msg.OfferType == "buy" {
			inc = 1
		}
		for _, calc := range pricesmonitor.AllCalculators {
			if calc.ID == id {
				calc.IncOffers(inc)
			}
		}
		fmt.Printf("New message: %s\n", msg.OfferType)
	}
}

//UpdateData is func
func UpdateData(id int, message OutcomingMessage) {
	for client := range clients {
		if clients[client] == id {
			err := client.WriteJSON(message)
			if err != nil {
				defer client.Close()
				delete(clients, client)
			}
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	id++
	file, _ := ioutil.ReadFile("./templates/graph-page.html")
	w.Write(file)
}

//StartServer is func
func StartServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./templates"))))

	router.HandleFunc("/", mhandler)
	router.HandleFunc("/graph{id:[0-9]+}", test)

	router.HandleFunc("/ws/graph{id:[0-9]+}", handleConnections)
	http.Handle("/", router)

	log.Println("starting http server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
