package netcommunicator

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//OutcomingMessage is a struct
type OutcomingMessage struct {
	Type        string  `json:"type"`
	OffersCount int     `json:"offersCount"`
	StockPrice  float64 `json:"stockPrice"`
}

//IncomingMessage is struct too
type IncomingMessage struct {
	OfferType string `json:"offerType"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan int)               // broadcast channel

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func replFunc(inp []byte) []byte {
	expr := regexp.MustCompile(`src="(.+?)"`)
	path := expr.FindStringSubmatch(string(inp))[1]
	rawData, err := ioutil.ReadFile("./templates/" + path)
	checkerr(err)
	return append([]byte("<script>\n"), append(rawData, []byte("\n</script>")...)...)
}

func insertScripts(pathToHTML string) string {
	rawData, err := ioutil.ReadFile("./templates/" + pathToHTML)
	checkerr(err)
	expr := regexp.MustCompile("<script type=\"text/javascript\" src=\".+?\"></script>")
	return string(expr.ReplaceAllFunc(rawData, replFunc))
}

func mhandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(insertScripts("index.html")))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	clients[ws] = true
	defer ws.Close()

	for {
		var msg IncomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		var inc int
		if msg.OfferType == "sell" {
			inc = -1
		} else if msg.OfferType == "buy" {
			inc = 1
		}
		broadcast <- inc
		fmt.Printf("New message: %s\n", msg.OfferType)
	}
}

func handleMessages() {
	curVal := 0
	for {
		curVal += <-broadcast

		for client := range clients {
			err := client.WriteJSON(OutcomingMessage{Type: "offers", OffersCount: curVal})
			if err != nil {
				log.Printf("error: %v", err)
				defer client.Close()
				delete(clients, client)
			}
		}
	}
}

func updateData() {
	value := 1.0
	for {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C

		value *= 1.1

		for client := range clients {
			err := client.WriteJSON(OutcomingMessage{Type: "gpoint", StockPrice: value})
			if err != nil {
				log.Printf("error: %v", err)
				defer client.Close()
				delete(clients, client)
			}
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("lol")
}

//StartServer is func
func StartServer() {
	insertScripts("index.html")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", mhandler)
	router.HandleFunc("/products/{id:[0-9]+}", test)

	http.HandleFunc("/ws", handleConnections)
	http.Handle("/", router)

	go updateData()
	go handleMessages()

	log.Println("starting http server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
