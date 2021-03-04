package database

import (
	"testing"

	"github.com/nebolax/LatenessStockExcahnge/general"
)

func TestAddStock(t *testing.T) {
	defer clearTest()

	var source = "testStock"

	err := AddStock("testStock", 0, 10)

	if !general.CheckError(err) {
		t.Error("Some error occurred: " + err.Error())
	}

	res, err := dataBase.Query("SELECT Name FROM stocks")

	if !general.CheckError(err) {
		t.Error("Some error occurred: " + err.Error())
	}

	for res.Next() {
		var data string
		_ = res.Scan(&data)

		if data == source {
			return
		}
	}

	t.Error("Name was not found")
}
