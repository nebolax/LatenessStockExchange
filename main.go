// Main package which unites all other parts
package main

import (
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/verifier"
)

// Main function to call all other functions
func main() {
	test()
	database.Init(database.StandartPath)
	verifier.CheckTest()
}
