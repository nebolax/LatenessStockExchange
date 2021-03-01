package database

import (
	"fmt"
	"time"
)

// Insert updated price to stock
func UpdatePrice(stockId int, price float64) error {
	if !initialized {
		return DatabaseError{"Database is not initialized"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO price_logs (stock_id, price, timestamp) values (%d, %f, %s)",
		stockId, price, time.Now().Format(datetimeFormat)))

	return err
}
