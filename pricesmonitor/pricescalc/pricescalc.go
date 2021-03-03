package pricescalc

import (
	"fmt"
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"math"
	"sync"
	"time"
)

const updateDatabasePeriod = 50

type stockDataHandler struct {
	CurStock      float64
	CurOffers     int
	SharesInTrade int
	mu            sync.Mutex
}

//RTPriceCalculator is a channel using which you will get stock prices in real time
type RTPriceCalculator struct {
	ID         int
	Name 	   string
	CurHandler *stockDataHandler
	LiveStock  chan float64
	LiveOffers chan int
	History    []float64
}

//IncOffers is a very good func
func (calc *RTPriceCalculator) IncOffers(val int) {
	calc.CurHandler.mu.Lock()
	calc.CurHandler.CurOffers += val
	calc.LiveOffers <- calc.CurHandler.CurOffers
	calc.CurHandler.mu.Unlock()
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
func CreatePriceCalculator(id int, startStock float64, offers, sharesInTrade int, name string) *RTPriceCalculator {
	curHandler := stockDataHandler{CurStock: startStock, CurOffers: offers, SharesInTrade: sharesInTrade}
	priceObs := RTPriceCalculator{id, name, &curHandler, make(chan float64), make(chan int), []float64{}}

	go updatePrices(&priceObs)

	return &priceObs
}

func updatePrices(observer *RTPriceCalculator) {
	leftTillUpdate := updateDatabasePeriod

	for {

		leftTillUpdate--

		if leftTillUpdate == 0 {
			fmt.Println(observer)
			fmt.Println("update started")

			leftTillUpdate = updateDatabasePeriod
			var price = observer.History[len(observer.History) - 1]
			fmt.Println(price)
			err := database.UpdatePrice(observer.ID, price)
			general.CheckError(err)

			fmt.Println("update finished")
		}

		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		observer.LiveStock <- observer.newPrice()
	}
}
