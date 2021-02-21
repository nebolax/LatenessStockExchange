package main

import (
	"github.com/nebolax/LatenessStockExcahnge/netcommunicator"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
)

func main() {
	for i := 0; i < len(pricesmonitor.AllCalculators); i++ {
		go procNewStocks(pricesmonitor.AllCalculators[i].LivePrice)
	}
	netcommunicator.StartServer()
}

func procNewStocks(ch chan float64) {
	newPrice := <-ch
	//TODO - save newPrice to database
	newPrice++
}
