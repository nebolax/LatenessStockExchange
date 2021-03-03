package database

import (
	"fmt"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"time"
)

// Id of gritscoin in table with stocks
const gritscoinId = 0

//Possible actions
const (
	sell = iota
	buy = iota
	dividends = iota
)

// A datetime format for database storage
const datetimeFormat = "2000-01-01 11:12:13"

// Count Amount of StockId stock owned by buyerId
func getResources(buyerId int, stockId int) (int, error) {
	count, err := dataBase.Query(fmt.Sprintf(
		"SELECT Amount FROM user_stock_ownerships WHERE (user_id = %d) AND (stock_id = %d)", buyerId, stockId))
	if !general.CheckError(err) {
		return 0, err
	}


	for count.Next() {
		var countInt int
		err = count.Scan(&countInt)
		if !general.CheckError(err) {
			return 0, err
		}

		return countInt, nil
	}

	return 0, nil

}

func getGritscoins(userId int) (float64, error) {
	count, err := dataBase.Query(fmt.Sprintf(
		"SELECT money FROM users WHERE (id = %d)", userId))
	if !general.CheckError(err) {
		return 0, err
	}

	for count.Next() {
		var amount float64
		err = count.Scan(&amount)
		if !general.CheckError(err) {
			return 0, err
		}
		return amount, nil
	}

	return 0, DatabaseError{"No such user"}
}

// Makes transaction between sellerId and buyerId, where
// sellerId sells Amount of StockId stocks to buyerId
func MakeTransaction(sellerId int, buyerId int, stockId int, amount int, currentPrice float64) error {
	totalDealPrice := currentPrice * float64(amount)

	buyerGritscoins, err := getGritscoins(buyerId)

	if !general.CheckError(err) {
		return err
	}

	buyerStocks, err := getResources(buyerId, stockId)

	if !general.CheckError(err) {
		return err
	}

	sellerGritscoins, err := getGritscoins(sellerId)

	if !general.CheckError(err) {
		return err
	}

	sellerStocks, err := getResources(sellerId, stockId)

	if !general.CheckError(err) {
		return err
	}

	if buyerGritscoins < totalDealPrice {
		return DatabaseError{"Buyer has not enough money for deal"}
	}

	if sellerStocks < amount {
		return DatabaseError{"Seller has not enough stocks for deal"}
	}

	baseStock := "UPDATE user_stock_ownerships SET Amount = %d WHERE (user_id = %d) AND (stock_id = %d)"

	baseMoney := "UPDATE users SET money = %f WHERE (user_id = %d)"

	updateSellerStocks := fmt.Sprintf(baseStock, sellerStocks - amount, sellerId, stockId)
	_, err = dataBase.Exec(updateSellerStocks)
	if !general.CheckError(err) {
		return err
	}

	updateBuyerCoins := fmt.Sprintf(baseMoney, buyerGritscoins - totalDealPrice, buyerId)
	_, err = dataBase.Exec(updateBuyerCoins)
	if !general.CheckError(err) {
		return err
	}

	updateSellerCoins := fmt.Sprintf(baseMoney, sellerGritscoins + totalDealPrice, sellerId)
	_, err = dataBase.Exec(updateSellerCoins)
	if !general.CheckError(err) {
		return err
	}

	updateBuyerStocks := fmt.Sprintf(baseStock, buyerStocks + amount, buyerId, stockId)
	_, err = dataBase.Exec(updateBuyerStocks)
	if !general.CheckError(err) {
		return err
	}

	addTransaction := "INSERT INTO transaction_logs (user_id, stock_id, Amount, money_spent, type, timestamp) " +
		"values (%d, %d, %d, %f, %d, %s)"
	_, err = dataBase.Exec(fmt.Sprintf(addTransaction, sellerId, stockId, amount,
		-1 * totalDealPrice, sell, time.Now().Format(datetimeFormat)))
	if !general.CheckError(err) {
		return err
	}

	_, err = dataBase.Exec(fmt.Sprintf(addTransaction, buyerId, stockId, amount,
		totalDealPrice, buy, time.Now().Format(datetimeFormat)))
	if !general.CheckError(err) {
		return err
	}

	return nil

}

// Gives dividends to all StockId owners
// Gives percentage (1.0 = 100%) of cost of stock to every owner for one stock
func Dividends(stockId int, percentage float64) error {
	owners, err := dataBase.Query(fmt.Sprintf(
		"SELECT amount, user_id FROM user_stock_ownerships WHERE (stock_id = %d)", stockId))

	if !general.CheckError(err) {
		return err
	}

	price, priceErr := dataBase.Query(fmt.Sprintf(
		"SELECT price FROM price_logs WHERE (id = %d) ORDER BY timestamp DESC LIMIT 1", stockId))
	if !general.CheckError(priceErr) {
		return priceErr
	}

	var priceData float64
	if price.Next() {
		priceErr = price.Scan(&priceData)
		if !general.CheckError(priceErr) {
			return priceErr
		}
	}

	for owners.Next() {
		var userId int
		var amount int

		err = owners.Scan(&amount, &userId)
		if !general.CheckError(err) {
			return err
		}

		countResponse, err := dataBase.Query(fmt.Sprintf(
			"SELECT money FROM user WHERE (user_id = %d)", userId))
		if !general.CheckError(err) {
			return err
		}

		if countResponse.Next() {
			var value float64
			err = countResponse.Scan(&value)
			if !general.CheckError(err) {
				return err
			}

			delta := float64(amount) * percentage * priceData
			value += delta

			_, err = dataBase.Exec(fmt.Sprintf(
				"UPDATE users SET Amount = %f WHERE (user_id = %d)", value, userId))

			if !general.CheckError(err) {
				return err
			}

			_, err = dataBase.Exec(fmt.Sprintf(
				"INSERT INTO transaction_logs (user_id, stock_id, Amount, money_spent, type, timestamp) " +
				"values (%d, %d, %d, %f, %d, %s)",
				userId, stockId, amount, -delta, dividends, time.Now().Format(datetimeFormat)))

			if !general.CheckError(err) {
				return err
			}
		} else {
			return DatabaseError{"This person has no gritscoin bill"}
		}
	}

	return nil
}
