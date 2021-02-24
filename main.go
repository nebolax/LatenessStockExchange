package main

import (
	"github.com/nebolax/LatenessStockExcahnge/netcommunicator"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

func main() {
	for i := 0; i < len(pricesmonitor.AllCalculators); i++ {
		go procNewStocks(pricesmonitor.AllCalculators[i])
	}
	netcommunicator.StartServer()
}

func procNewStocks(calc *pricescalc.RTPriceCalculator) {
	for {
		newPrice := <-calc.LivePrice
		//TODO - save newPrice to database
		newPrice++
	}
}
