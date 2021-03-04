package pricescalc

import (
	"math"
	"sync"
	"time"
)

type stockDataHandler struct {
	CurStock      float64
	CurOffers     int
	SharesInTrade int
	mu            sync.Mutex
}

//RTPriceCalculator is a channel using which you will get stock prices in real time
type RTPriceCalculator struct {
	ID         int
	CurHandler *stockDataHandler
	LiveStock  chan float64
	LiveOffers chan int
	History    []float64
	offers     map[int]int
}

//ReqOffer is a very good func
func (calc *RTPriceCalculator) ReqOffer(askerID, amount int) int {
	calc.CurHandler.mu.Lock()
	defer calc.CurHandler.mu.Unlock()
	if _, ok := calc.offers[askerID]; ok {
		calc.offers[askerID] += amount
	} else {
		calc.offers[askerID] = amount
	}
	if calc.offers[askerID] == 0 {
		delete(calc.offers, askerID)
		return 0
	}
	calc.CurHandler.CurOffers += amount
	calc.LiveOffers <- calc.CurHandler.CurOffers

	return calc.offers[askerID]
}

//PersonOffers is a func
func (calc RTPriceCalculator) PersonOffers(personID int) int {
	if _, ok := calc.offers[personID]; ok {
		return calc.offers[personID]
	} else {
		return 0
	}
}

func (calc *RTPriceCalculator) newPrice() float64 {
	dh := calc.CurHandler
	speed := float64(dh.CurOffers) * (1.0 - math.Pow(0.5, 1.0/900))
	dh.CurStock *= (1 + speed)
	calc.History = append(calc.History, dh.CurStock)

	if len(calc.History) > 20 {
		calc.History = calc.History[len(calc.History)-20:]
	}
	return dh.CurStock
}

//CreatePriceCalculator is func
func CreatePriceCalculator(id int, startStock float64, offers, sharesInTrade int) *RTPriceCalculator {
	curHandler := stockDataHandler{CurStock: startStock, CurOffers: offers, SharesInTrade: sharesInTrade}
	priceObs := RTPriceCalculator{id, &curHandler, make(chan float64), make(chan int), []float64{}, make(map[int]int)}

	go updatePrices(&priceObs)

	return &priceObs
}

func updatePrices(observer *RTPriceCalculator) {
	for {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		observer.LiveStock <- observer.newPrice()
	}
}
