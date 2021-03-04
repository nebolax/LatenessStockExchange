package database

import (
	"fmt"

	"github.com/nebolax/LatenessStockExcahnge/general"
)

// Register new type of stock
func AddStock(name string, traderId int, totalCount int, startCost float64) error {
	if !initialized {
		return DatabaseError{"Database is not initialized!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO stocks (name, user_id, total_count) values ('%s', %d, %d)",
		name, traderId, totalCount))

	idReq, err := dataBase.Query(fmt.Sprintf(
		"SELECT id FROM stocks WHERE (name = '%s')", name))

	if !general.CheckError(err) {
		return err
	}

	defer idReq.Close()

	if idReq.Next() {
		var id int
		err = idReq.Scan(&id)

		if !general.CheckError(err) {
			return err
		}

		err = UpdatePrice(id, startCost)

		return err
	} else {
		return DatabaseError{"Some cringy cringe!!! There is no just created stock"}
	}
}

func GetAllStocks() ([]OneStockInfo, error) {
	if !initialized {
		return nil, DatabaseError{"Database is not initialized"}
	}

	result, err := dataBase.Query("SELECT id, name, total_count FROM stocks")

	if !general.CheckError(err) {
		return nil, err
	}
	defer result.Close()

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
