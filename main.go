// Main package which unites all other parts
package main

import (
	NetComms "github.com/nebolax/LatenessStockExcahnge/netcommunicator"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

// Main function to call all other functions
func main() {
	// for i := 0; i < len(pricesmonitor.AllCalculators); i++ {
	// 	procNewData(pricesmonitor.AllCalculators[i])
	// }
	// NetComms.StartServer()
	send()
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
