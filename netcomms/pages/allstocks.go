package pages

import (
	"fmt"
	"html/template"
	"math"
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/pricescalc"
)

//AllStocksObserver func
func AllStocksObserver(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/all-stocks-observer.html")
	var ar []database.OneStockInfo
	fmt.Println("len: ", len(pricescalc.AllCalculators))
	for calc := range pricescalc.AllCalculators {
		ar = append(ar, database.OneStockInfo{ID: calc.ID, Name: calc.Name, CurPrice: math.Round(calc.CurHandler.CurStock*100) / 100})
	}
	tmpl.Execute(w, ar)
}
