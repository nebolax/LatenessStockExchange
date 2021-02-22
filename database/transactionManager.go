package database

import (
	"fmt"
	"time"
)

const gritscoinId = 0

const (
	sell = iota
	buy = iota
	dividends = iota
)

const datetimeFormat = "2000-01-01 11:12:13"

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

func MakeTransaction(sellerId int, buyerId int, stockId int, amount int) error {
	price, priceErr := dataBase.Query(fmt.Sprintf(
		"SELECT price FROM price_logs WHERE (id = %d) ORDER BY timestamp DESC LIMIT 1", stockId))
	if !checkError(priceErr) {
		return priceErr
	}

	if price.Next() {
		var priceData float64
		price.Scan(&priceData)

		totalDealPrice := int(priceData* float64(amount)) + 1

		buyerGritscoins := getResources(buyerId, gritscoinId)
		buyerStocks := getResources(buyerId, stockId)

		sellerGritscoins := getResources(sellerId, gritscoinId)
		sellerStocks := getResources(sellerId, stockId)

		if buyerGritscoins < totalDealPrice {
			return DatabaseError{"Buyer has not enough money for deal"}
		}

		if sellerStocks < amount {
			return DatabaseError{"Seller has not enough stocks for deal"}
		}

		base := "UPDATE user_stock_ownerships SET amount = %d WHERE (user_id = %d) AND (stock_id = %d)"

		updateSellerStocks := fmt.Sprintf(base, sellerStocks - amount, sellerId, stockId)
		_, err := dataBase.Exec(updateSellerStocks)
		if !checkError(err) {
			return err
		}

		updateBuyerCoins := fmt.Sprintf(base, buyerGritscoins - totalDealPrice, buyerId, gritscoinId)
		_, err = dataBase.Exec(updateBuyerCoins)
		if !checkError(err) {
			return err
		}

		updateSellerCoins := fmt.Sprintf(base, sellerGritscoins + totalDealPrice, sellerId, gritscoinId)
		_, err = dataBase.Exec(updateSellerCoins)
		if !checkError(err) {
			return err
		}

		updateBuyerStocks := fmt.Sprintf(base, buyerStocks + amount, buyerId, stockId)
		_, err = dataBase.Exec(updateBuyerStocks)
		if !checkError(err) {
			return err
		}

		addTransaction := "INSERT INTO transaction_logs (user_id, stock_id, amount, money_spent, type, timestamp) " +
			"values (%d, %d, %d, %d, %d, %s)"
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
	} else {
		return DatabaseError{"There is no such stockId"}
	}
}

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
			"SELECT amount FROM user_stock_ownerships WHERE (user_id = %d) AND (stock_id = %d)",
			userId, gritscoinId))
		if !checkError(err) {
			return err
		}

		if countResponse.Next() {
			var value int
			err = countResponse.Scan(&value)
			if !checkError(err) {
				return err
			}

			delta := int(float64(amount) * percentage * priceData) + 1
			value += delta

			_, err = dataBase.Exec(fmt.Sprintf(
				"UPDATE user_stock_ownerships SET amount = %d WHERE (user_id = %d) AND (stock_id = %d)",
				value, userId, gritscoinId))

			if !checkError(err) {
				return err
			}
		} else {
			return DatabaseError{"This person has no gritscoin bill"}
		}
	}

	return nil
}
