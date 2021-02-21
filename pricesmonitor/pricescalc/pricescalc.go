package pricescalc

import (
	"math"
	"sync"
	"time"
)

type stockDataHandler struct {
	curStock      float64
	curOffers     int
	sharesInTrade int
	mu            sync.Mutex
}

//RTPriceCalculator is a channel using which you will get stock prices in real time
type RTPriceCalculator struct {
	Name       int
	curHandler *stockDataHandler
	LivePrice  chan float64
}

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

//CreatePriceCalculator is func
func CreatePriceCalculator(id int, startStock float64, offers, sharesInTrade int) RTPriceCalculator {
	curHandler := stockDataHandler{curStock: startStock, curOffers: offers, sharesInTrade: sharesInTrade}
	livePrice := make(chan float64)
	priceObs := RTPriceCalculator{id, &curHandler, livePrice}

	go updatePrices(&priceObs)

	return priceObs
}

func updatePrices(observer *RTPriceCalculator) {
	for {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		observer.LivePrice <- observer.curHandler.newPrice()
	}
}
