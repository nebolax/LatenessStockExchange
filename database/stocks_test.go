package database

import (
	"testing"
)

func TestAddStock(t *testing.T) {
	defer clearTest()
	Init(testPath)

	var source = "testStock"

	err := AddStock("testStock", 0)

	if !checkError(err) {
		t.Error("Some error occurred: " + err.Error())
	}

	res, err := dataBase.Query("SELECT name FROM stocks")

	if !checkError(err) {
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
