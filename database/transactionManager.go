package database

import (
	"fmt"
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

// Count amount of stockId stock owned by buyerId
func getResources(buyerId int, stockId int) int {
	count, err := dataBase.Query(fmt.Sprintf(
		"SELECT amount from user_stock_ownership WHERE (user_id = %d) AND (stock_id = %d)", buyerId, stockId))
	if !checkError(err) {
		return 0
	}

	if count.Next() {
		var countInt int
		count.Scan(&countInt)
		return countInt
	} else {
		return 0
	}
}

func getGritscoins(userId int) float64 {
	count, err := dataBase.Query(fmt.Sprintf(
		"SELECT money from users WHERE (user_id = %d)", userId))
	if !checkError(err) {
		return 0
	}

	if count.Next() {
		var amount float64
		count.Scan(&amount)
		return amount
	} else {
		return 0
	}
}

// Makes transaction between sellerId and buyerId, where
// sellerId sells amount of stockId stocks to buyerId
func MakeTransaction(sellerId int, buyerId int, stockId int, amount int, currentPrice float64) error {
	totalDealPrice := currentPrice * float64(amount)

	buyerGritscoins := getGritscoins(buyerId)
	buyerStocks := getResources(buyerId, stockId)

	sellerGritscoins := getGritscoins(sellerId)
	sellerStocks := getResources(sellerId, stockId)

	if buyerGritscoins < totalDealPrice {
		return DatabaseError{"Buyer has not enough money for deal"}
	}

	if sellerStocks < amount {
		return DatabaseError{"Seller has not enough stocks for deal"}
	}

	baseStock := "UPDATE user_stock_ownerships SET amount = %d WHERE (user_id = %d) AND (stock_id = %d)"

	baseMoney := "UPDATE users SET money = %f WHERE (user_id = %d)"

	updateSellerStocks := fmt.Sprintf(baseStock, sellerStocks - amount, sellerId, stockId)
	_, err := dataBase.Exec(updateSellerStocks)
	if !checkError(err) {
		return err
	}

	updateBuyerCoins := fmt.Sprintf(baseMoney, buyerGritscoins - totalDealPrice, buyerId)
	_, err = dataBase.Exec(updateBuyerCoins)
	if !checkError(err) {
		return err
	}

	updateSellerCoins := fmt.Sprintf(baseMoney, sellerGritscoins + totalDealPrice, sellerId)
	_, err = dataBase.Exec(updateSellerCoins)
	if !checkError(err) {
		return err
	}

	updateBuyerStocks := fmt.Sprintf(baseStock, buyerStocks + amount, buyerId, stockId)
	_, err = dataBase.Exec(updateBuyerStocks)
	if !checkError(err) {
		return err
	}

	addTransaction := "INSERT INTO transaction_logs (user_id, stock_id, amount, money_spent, type, timestamp) " +
		"values (%d, %d, %d, %f, %d, %s)"
	_, err = dataBase.Exec(fmt.Sprintf(addTransaction, sellerId, stockId, amount,
		-1 * totalDealPrice, sell, time.Now().Format(datetimeFormat)))
	if !checkError(err) {
		return err
	}

	_, err = dataBase.Exec(fmt.Sprintf(addTransaction, buyerId, stockId, amount,
		totalDealPrice, buy, time.Now().Format(datetimeFormat)))
	if !checkError(err) {
		return err
	}

	return nil

}

// Gives dividends to all stockId owners
// Gives percentage (1.0 = 100%) of cost of stock to every owner for one stock
func Dividends(stockId int, percentage float64) error {
	owners, err := dataBase.Query(fmt.Sprintf(
		"SELECT (amount, user_id) FROM user_stock_ownerships WHERE (stock_id = %d)", stockId))

	if !checkError(err) {
		return err
	}

	price, priceErr := dataBase.Query(fmt.Sprintf(
		"SELECT price FROM price_logs WHERE (id = %d) ORDER BY timestamp DESC LIMIT 1", stockId))
	if !checkError(priceErr) {
		return priceErr
	}

	var priceData float64
	if price.Next() {
		priceErr = price.Scan(&priceData)
		if !checkError(priceErr) {
			return priceErr
		}
	}

	for owners.Next() {
		var userId int
		var amount int

		err = owners.Scan(&amount, &userId)
		if !checkError(err) {
			return err
		}

		countResponse, err := dataBase.Query(fmt.Sprintf(
			"SELECT money FROM user WHERE (user_id = %d)", userId))
		if !checkError(err) {
			return err
		}

		if countResponse.Next() {
			var value float64
			err = countResponse.Scan(&value)
			if !checkError(err) {
				return err
			}

			delta := float64(amount) * percentage * priceData
			value += delta

			_, err = dataBase.Exec(fmt.Sprintf(
				"UPDATE users SET amount = %f WHERE (user_id = %d)", value, userId))

			if !checkError(err) {
				return err
			}

			_, err = dataBase.Exec(fmt.Sprintf(
				"INSERT INTO transaction_logs (user_id, stock_id, amount, money_spent, type, timestamp) " +
				"values (%d, %d, %d, %f, %d, %s)",
				userId, stockId, amount, -delta, dividends, time.Now().Format(datetimeFormat)))

			if !checkError(err) {
				return err
			}
		} else {
			return DatabaseError{"This person has no gritscoin bill"}
		}
	}

	return nil
}
