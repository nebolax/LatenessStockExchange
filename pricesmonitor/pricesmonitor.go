package pricesmonitor

import (
	"math/rand"
	"time"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"

	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor/pricescalc"
)

//AllCalculators contains all prices calculators
var AllCalculators []*pricescalc.RTPriceCalculator

func Init() {
	rand.Seed(time.Now().UnixNano())

	AllCalculators = make([]*pricescalc.RTPriceCalculator, 0)

	allStocks, err := database.GetAllStocks()

	if !general.CheckError(err) {
		return
	}

	for _, stock := range allStocks {
		newCalc := pricescalc.CreatePriceCalculator(
			stock.ID, stock.CurPrice, 0, stock.TotalCount, stock.Name)

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
