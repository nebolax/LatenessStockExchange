package pricesmonitor

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

var allCalculators []pricescalc.RTPriceCalculator

func init() {
	for i := 0; i < 10; i++ {
		mshares := rand.Int()%20 - 10
		newCalc := pricescalc.CreatePriceCalculator(strconv.Itoa(i), rand.Float64()*100, mshares, int(math.Abs(float64(mshares))))
		allCalculators = append(allCalculators, newCalc)
	}

	for i := 0; i < len(allCalculators); i++ {
		go observeStock(&allCalculators[i])
	}
}

func observeStock(calc *pricescalc.RTPriceCalculator) {
	for {
		nprice := <-calc.LivePrice
		fmt.Printf("Stock '%s', price: %f", calc.Name, nprice)
	}
}
