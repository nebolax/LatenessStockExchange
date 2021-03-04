package procs

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs/connsinc"
	"github.com/nebolax/LatenessStockExcahnge/pricescalc"
)

type sockClient struct {
	userID int
	calc   *pricescalc.RTPriceCalculator
	sock   *websocket.Conn
	mu     sync.Mutex
}

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
	Type         string    `json:"type"`
	History      []float64 `json:"history"`
	PublicOffers int       `json:"publicOffers"`
	PersonOffers int       `json:"personOffers"`
}

var clients = make(map[int]*sockClient)

//SendtoGraphObservers is func
func SendtoGraphObservers(graphID int, message interface{}) {
	for id, client := range clients {
		if client.calc.ID == graphID {
			writeSingleMessage(id, message)
		}
	}
}

//sendtoUserDevices is func
func sendtoUserDevices(userID int, message interface{}) {
	for connID, client := range clients {
		if client.userID == userID {
			writeSingleMessage(connID, message)
		}
	}
}

//writeSingleMessage is func
func writeSingleMessage(connID int, message interface{}) {
	client := clients[connID]
	client.mu.Lock()
	err := client.sock.WriteJSON(message)
	if err != nil {
		delClient(connID)
	}
	client.mu.Unlock()
}

func delClient(connID int) {
	clients[connID].sock.Close()
	delete(clients, connID)
}

//ReadSingleMessage is func
func readSingleMessage(connID int) (incomingMessage, bool) {
	var msg incomingMessage
	err := clients[connID].sock.ReadJSON(&msg)
	if err != nil {
		delClient(connID)
		return incomingMessage{}, false
	}
	return msg, true
}

func procIncomingMessages(connID int) {
	for {
		msg, ok := readSingleMessage(connID)
		if ok {
			var offs int
			if msg.OfferType == "sell" {
				offs = clients[connID].calc.ReqOffer(connID, -1)
			} else if msg.OfferType == "buy" {
				offs = clients[connID].calc.ReqOffer(connID, 1)
			}

			sendtoUserDevices(clients[connID].userID, OutcomingMessage{Type: "personOffers", OffersCount: offs})
		}
	}
}

//RegNewGraphObserver is func
func RegNewGraphObserver(ws *websocket.Conn, calc *pricescalc.RTPriceCalculator, userID int) {
	connID := connsinc.NewID()
	clients[connID] = &sockClient{userID: userID, calc: calc, sock: ws}
	writeSingleMessage(connID, graphPageSetup{"setup", calc.History, calc.CurHandler.CurOffers, calc.PersonOffers(userID)})
	go procIncomingMessages(connID)
}
