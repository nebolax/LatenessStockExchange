package netcommunicator

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//Message is a struct
type Message struct {
	Type        string  `json:"type"`
	OffersCount int     `json:"offersCount"`
	StockPrice  float64 `json:"stockPrice"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

func mhandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		println("error", err.Error())
	} else {
		tmpl.Execute(w, nil)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	clients[ws] = true
	println("accepted new client")
}

// func handleMessages() {
// 	for {
// 		msg := <-broadcast

// 		for client := range clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				defer client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }

func updateData() {
	value := 1.0
	for {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C

		value *= 1.1

		for client := range clients {
			err := client.WriteJSON(Message{Type: "gpoint", StockPrice: value, OffersCount: int(value)})
			if err != nil {
				log.Printf("error: %v", err)
				defer client.Close()
				delete(clients, client)
			} else {
			}
		}
		fmt.Printf("Sent value %f; Current clients count: %d\n", value, len(clients))
	}
}

//StartServer is func
func StartServer() {
	fs := http.FileServer(http.Dir("templates"))
	//http.HandleFunc("/", mhandler)
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go updateData()

	log.Println("starting http server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
