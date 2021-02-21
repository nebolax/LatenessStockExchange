package main

import (
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/verifier"
)

func main() {
	test()
	database.Init()
	verifier.CheckTest()
}
