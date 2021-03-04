package database

import (
	"math"
	"testing"

	"github.com/nebolax/LatenessStockExcahnge/general"
)

func TestUpdatePrice(t *testing.T) {
	defer clearTest()

	var data = 179.179

	err := UpdatePrice(1, data)

	if err != nil {
		t.Error("Some error occurred: " + err.Error())
	}

	value, err := GetLatestPrice(1)

	if !general.CheckError(err) {
		t.Error("Some error occurred: " + err.Error())
	}

	if math.Abs(data-value) > 0.01 {
		t.Error("Source value not equals result value")
	}
}
