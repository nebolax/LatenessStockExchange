// Main package which unites all other parts
package main

import (
	"github.com/nebolax/LatenessStockExcahnge/database"
	netcomms "github.com/nebolax/LatenessStockExcahnge/netcomms"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
)

// Main function to call all other functions
func main() {
	database.Init(database.StandartPath)
	pricesmonitor.Init()
	netcomms.StartServer()
}
