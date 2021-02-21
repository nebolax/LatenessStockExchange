package main

import (
	"github.com/nebolax/LatenessStockExcahnge/Database"
	"github.com/nebolax/LatenessStockExcahnge/Verifier"
)

func main() {
	test()
	Database.TestDb()
	Verifier.CheckTest()
}