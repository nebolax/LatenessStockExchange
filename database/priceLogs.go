package database

import (
	"fmt"
	"time"

	"github.com/nebolax/LatenessStockExcahnge/general"
)

// Insert updated price to stock
func UpdatePrice(stockId int, price float64) error {
	if !initialized {
		return DatabaseError{"Database is not initialized"}
	}

	query := fmt.Sprintf(
		"INSERT INTO price_logs (stock_id, price, timestamp) VALUES (%d, %f, '%s')",
		stockId, price, general.TimeToString(time.Now()))

	_, err := dataBase.Exec(query)

	return err
}

func GetLatestPrice(stockId int) (float64, error) {
	res, err := dataBase.Query(
		fmt.Sprintf("SELECT price FROM price_logs WHERE "+
			"(timestamp = (SELECT MAX(timestamp) FROM price_logs WHERE (stock_id = %d))) "+
			"AND (stock_id = %d)", stockId, stockId))

	if !general.CheckError(err) {
		return 0, err
	}

	defer res.Close()

	if res.Next() {
		var value float64
		err = res.Scan(&value)

		if !general.CheckError(err) {
			return 0, err
		}

		return value, nil
	} else {
		return 0, DatabaseError{"There are no records about such StockId"}
	}
}
