package main

import (
	"github.com/nebolax/LatenessStockExcahnge/pricescalc"
)

func main() {
	go monitorPrice()
	for {

	}
}

func monitorPrice() {
	for price := range pricescalc.LivePrice {
		println(price)
	}
}
