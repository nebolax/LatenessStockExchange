package pricescalc

import (
	"math"
	"sync"
	"time"
)

//LivePrice is a channel using which you will get stock prices in real time
var LivePrice chan float64

type stockDataHandler struct {
	curStock      float64
	mu            sync.Mutex
	curOffers     int
	sharesInTrade int
}

var curHandler stockDataHandler

func (dh *stockDataHandler) incOffers(val int) {
	dh.mu.Lock()
	dh.curOffers += val
	dh.mu.Unlock()
}

func (dh *stockDataHandler) newPrice() float64 {
	speed := float64(dh.curOffers) * (1.0 - math.Pow(0.5, 1.0/900))
	dh.curStock *= (1 + speed)
	return dh.curStock
}

func init() {
	curHandler = stockDataHandler{curStock: 5.0, curOffers: 10, sharesInTrade: 30}
	LivePrice = make(chan float64)
	go updatePrices()
}

func updatePrices() {
	for {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		LivePrice <- curHandler.newPrice()
	}
}
