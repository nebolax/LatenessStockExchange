package pricesmonitor

import (
	"math"
	"math/rand"
	"time"

	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

//AllCalculators contains all prices calculators
var AllCalculators []*pricescalc.RTPriceCalculator

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		mshares := rand.Int()%20 - 10
		newCalc := pricescalc.CreatePriceCalculator(i, rand.Float64()*100, mshares, int(math.Abs(float64(mshares))))
		AllCalculators = append(AllCalculators, newCalc)
	}
}

//CalcByID id a func
func CalcByID(id int) *pricescalc.RTPriceCalculator {
	for _, calc := range AllCalculators {
		if calc.ID == id {
			return calc
		}
	}

	return nil
}
