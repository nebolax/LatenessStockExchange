package database

import (
	"fmt"
	"github.com/nebolax/LatenessStockExcahnge/general"
)

// Register new type of stock
func AddStock(name string, traderId int, totalCount int) error {
	if !initialized {
		return DatabaseError{"Database is not initialized!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO stocks (name, user_id, total_count) values ('%s', %d, %d)",
		name, traderId, totalCount))

	return err
}

func GetAllStocks() ([]OneStockInfo, error) {
	if !initialized {
		return nil, DatabaseError{"Database is not initialized"}
	}

	result, err := dataBase.Query("SELECT id, name, total_count FROM stocks")

	if !general.CheckError(err) {
		return nil, err
	}

	var allStocks = make([]OneStockInfo, 0)

	for result.Next() {
		var res = OneStockInfo{}
		err = result.Scan(&res.ID, &res.Name, &res.TotalCount)

		if !general.CheckError(err) {
			return nil, err
		}

		res.CurPrice, err = GetLatestPrice(res.ID)

		if !general.CheckError(err) {
			return nil, err
		}

		allStocks = append(allStocks, res)
	}

	return allStocks, nil
}
