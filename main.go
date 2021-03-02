// Main package which unites all other parts
package main

import (
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/verifier"
	NetComms "github.com/nebolax/LatenessStockExcahnge/netcommunicator"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

// Main function to call all other functions
func main() {
	database.Init(database.StandartPath)
	verifier.CheckTest()
	for i := 0; i < len(pricesmonitor.AllCalculators); i++ {
		procNewData(pricesmonitor.AllCalculators[i])
	}
	NetComms.StartServer()
}

func procNewData(calc *pricescalc.RTPriceCalculator) {
	go procNewStocks(calc.ID, calc.LiveStock)
	go procNewOffers(calc.ID, calc.LiveOffers)
}

func procNewStocks(id int, ch chan float64) {
	for {
		newStock := <-ch
		NetComms.UpdateData(id, NetComms.OutcomingMessage{Type: "gpoint", StockPrice: newStock})
	}
}

func procNewOffers(id int, ch chan int) {
	for {
		newOffers := <-ch
		NetComms.UpdateData(id, NetComms.OutcomingMessage{Type: "offers", OffersCount: newOffers})
	}
}
