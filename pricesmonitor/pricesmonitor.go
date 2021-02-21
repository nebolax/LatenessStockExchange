package pricesmonitor

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

//AllCalculators contains all prices calculators
var AllCalculators []*pricescalc.RTPriceCalculator

func init() {
	for i := 0; i < 10; i++ {
		mshares := rand.Int()%20 - 10
		newCalc := pricescalc.CreatePriceCalculator(i, rand.Float64()*100, mshares, int(math.Abs(float64(mshares))))
		AllCalculators = append(AllCalculators, &newCalc)
	}

	// for i := 0; i < len(allCalculators); i++ {
	// 	go observeStock(&allCalculators[i])
	// }
}

func observeStock(calc *pricescalc.RTPriceCalculator) {
	for {
		nprice := <-calc.LivePrice
		fmt.Printf("Stock '%d', price: %f", calc.Name, nprice)
	}
}
