package pricesmonitor

import (
	"math/rand"
	"time"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs"
	"github.com/nebolax/LatenessStockExcahnge/pricescalc"
)

//Init function
func Init() {
	rand.Seed(time.Now().UnixNano())

	allStocks, err := database.GetAllStocks()

	if !general.CheckError(err) {
		return
	}

	for _, stock := range allStocks {
		newCalc := pricescalc.CreatePriceCalculator(
			stock.ID, stock.CurPrice, 0, stock.TotalCount, stock.Name)
		procNewData(newCalc)
		pricescalc.AllCalculators[newCalc] = true
	}
}

func procNewData(calc *pricescalc.RTPriceCalculator) {
	go procNewStocks(calc.ID, calc.LiveStock)
	go procNewOffers(calc.ID, calc.LiveOffers)
}

func procNewStocks(graphID int, ch chan float64) {
	for {
		newStock := <-ch
		procs.SendtoGraphObservers(graphID, procs.OutcomingMessage{Type: "gpoint", StockPrice: newStock})
	}
}

func procNewOffers(graphID int, ch chan int) {
	for {
		newOffers := <-ch
		procs.SendtoGraphObservers(graphID, procs.OutcomingMessage{Type: "publicOffers", OffersCount: newOffers})
	}
}

//CalcByID id a func
func CalcByID(id int) *pricescalc.RTPriceCalculator {
	for calc := range pricescalc.AllCalculators {
		if calc.ID == id {
			return calc
		}
	}

	return nil
}
