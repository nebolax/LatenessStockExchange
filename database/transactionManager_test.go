package database

import (
	"fmt"
	"math"
	"testing"
)

func TestGetResources(t *testing.T) {
	defer clearTest()
	Init(testPath)
	userId := 123
	stockId := 234
	amount := 222
	_, _ = dataBase.Exec(fmt.Sprintf(
		"INSERT INTO user_stock_ownerships (user_id, stock_id, amount) " +
			"values (%d, %d, %d)", userId, stockId, amount))

	result, err := getResources(userId, stockId)

	if !checkError(err) {
		t.Error(err.Error())
	}

	if result != amount {
		t.Error(fmt.Sprintf("result and source do not match: %d / %d", result, amount))
	}
}

func TestGetGritscoins(t *testing.T) {
	defer clearTest()
	Init(testPath)
	userId := 123
	amount := 222.333
	_, _ = dataBase.Exec(fmt.Sprintf(
		"INSERT INTO users " +
			"(id, username, email, password_hash, password_salt, money) " +
			"values (%d, 'pes', 'pes@pes.pes', '!@#$1111', 'qwerty', %f)", userId, amount))

	result, err := getGritscoins(userId)

	if !checkError(err) {
		t.Error(err.Error())
	}

	if result != amount {
		t.Error(fmt.Sprintf("result and source do not match: %f / %f", result, amount))
	}
}

func prepareForTestMakeTransaction(){
	Init(testPath)

	_, _ = dataBase.Exec("INSERT INTO users " +
		"(id, username, email, password_hash, password_salt, money) " +
		"values (1, 'pes1', 'pes1@pes.pes', '!@#$1111', 'qwerty', 100.0)")

	_, _ = dataBase.Exec("INSERT INTO users " +
		"(id, username, email, password_hash, password_salt, money) " +
		"values (2, 'pes2', 'pes2@pes.pes', '!@#$1111', 'qwerty', 110.0)")


	_, _ = dataBase.Exec(
		"INSERT INTO user_stock_ownerships (user_id, stock_id, amount) values (1, 1, 10)")
}


func TestMakeTransactionFailMoney(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()

	var price = 553.5

	err := MakeTransaction(1, 2, 1, 2, price)

	if checkError(err) {
		t.Error("No error!!!")
	} else{
		if err.Error() != "Buyer has not enough money for deal" {
			t.Error("Wrong error message: " + err.Error())
		}
	}
}

func TestMakeTransactionFailDBStock(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()

	var price = 53.5

	err := MakeTransaction(1, 2, 179, 2, price)

	if checkError(err) {
		t.Error("No error occurred!!")
	}
}

func TestMakeTransactionFailDBUser(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()

	var price = 53.5

	err := MakeTransaction(179, 2, 1, 2, price)

	if checkError(err) {
		t.Error("No error occurred!!")
	}
}

func TestMakeTransactionFailStocks(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()
	var price = 1.5

	err := MakeTransaction(1, 2, 1, 20, price)

	if checkError(err) {
		t.Error("No error occurred!!")
	}
}

func TestMakeTransactionOk(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()

	var price = 20.5
	err := MakeTransaction(1, 2, 1, 2, price)

	if !checkError(err) {
		t.Error(err.Error())
	}
}

func TestDividends(t *testing.T) {
	defer clearTest()
	prepareForTestMakeTransaction()

	var price = 123.1

	err := UpdatePrice(1, price)

	if !checkError(err) {
		t.Error(err.Error())
	}

	var percentage = 0.5

	err = Dividends(1, percentage)

	if !checkError(err) {
		t.Error(err.Error())
	}

	res, err := dataBase.Query("SELECT (money) FROM users WHERE (id = 1)")

	if !checkError(err) {
		t.Error(err.Error())
	}

	if res.Next() {
		var money float64
		_ = res.Scan(&money)

		if math.Abs((money - 100.0) - price * 10 * percentage) > 0.01 {
			t.Error("Dividends does not match!!")
		}
	} else {
		t.Error("No user in db")
	}
}
