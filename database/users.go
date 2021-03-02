package database

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

// Name of batadase table with users
const tableName = "users"
// Key which is used in hash function
var hashKey = "YaPёSiyPЁS"
// The length of salt which would be added to password and hashKey
const saltLength = 16
// Able characters to use in salt
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate new random object with seed of current time
var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Generate random string of given length of given charset (copypasted from stackoverflow)
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Generate random string of given length (copypasted from stackoverflow)
func randomString(length int) string {
	return stringWithCharset(length, charset)
}

// Some email checking function (mb it's better to make it more intellectual)
func checkEmail(email string) bool {
	var template = regexp.MustCompile("[*]+@179.ru")
	return template.Match([]byte(email))
}

// Adds user with given email and password and nickname to batadase
func AddUser(nickname string, email string, password string) error{
	if !checkEmail(email) {
		return DatabaseError{"Email address is incorrect!"}
	}

	salt := randomString(saltLength)
	resultPassword := password + hashKey + salt
	hasher := sha512.New()
	hasher.Write([]byte(resultPassword))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if !initialized {
		return DatabaseError{"Database is not initialized yet!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO users (username, email, password_hash, password_salt, money) values ('%s', '%s', '%s', '%s', 0);",
		nickname, email, hash, salt))

	return err
}

// Tries to login user from sql query using given password. Returns error if something went wrong
// Returns nil if everything is OK
func login(user *sql.Rows, password string) error {
	if user.Next() {
		var salt string
		var passwordHash string

		err := user.Scan(&salt, &passwordHash)

		if !checkError(err) {
			return err
		}

		resultString := password + hashKey + salt

		hasher := sha512.New()
		hasher.Write([]byte(resultString))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		if hash != passwordHash {
			return DatabaseError{"Password does not match!"}
		} else {
			return nil
		}

	} else {
		return DatabaseError{"There is no such user"}
	}
}

// Tries to login user with given nickname by given password
// Returns error if something went wrong
// Returns nil if everything is OK
func LoginByNickname(nickname string, password string) error {
	if !initialized {
		return DatabaseError{"Database is not initialized"}
	}

	query, err := dataBase.Query(fmt.Sprintf("SELECT (password_salt, password_hash) FROM users WHERE (username = '%s')", nickname))

	if !checkError(err) {
		return err
	}

	return login(query, password)
}

// Tries to login user with given email by given password
// Returns error if something went wrong
// Returns nil if everything is OK
func LoginByEmail(email string, password string) error {
	if !initialized {
		return DatabaseError{"Database is not initialized"}
	}

	query, err := dataBase.Query(fmt.Sprintf("SELECT (password_salt, password_hash) FROM users WHERE (email = '%s')", email))

	if !checkError(err) {
		return err
	}

	return login(query, password)
}

// Returns map from stockId to amount - ownerships of userId
// If error occurred, returns nil as map and error as second object
func GetInvestmentPortfolio(userId int) (map[int]int, error) {
	if !initialized {
		return nil, DatabaseError{"Database is not initialized"}
	}

	query, err := dataBase.Query(fmt.Sprintf(
		"SELECT (stock_id, amount) FROM user_stock_ownerships WHERE (user_id = %d)", userId))

	if !checkError(err) {
		return nil, err
	}

	var result = make(map[int]int)

	for query.Next() {
		var stockId int
		var amount int

		err = query.Scan(&stockId, &amount)

		if !checkError(err) {
			return nil, err
		}

		result[stockId] = amount
	}

	return result, nil

}

// Returns more info about investment portfolio of userId (see Ownership for more information)
func GetInvestmentPortfolioPretty(userId int) ([]Ownership, error){
	data, err := GetInvestmentPortfolio(userId)
	if !checkError(err) {
		return nil, err
	}

	var result = make([]Ownership, 0)

	for key, value := range data {
		var ownership = Ownership{key, "", value, 0}

		res, err := dataBase.Query(fmt.Sprintf("SELECT name FROM stocks WHERE (id = %d)", key))
		if !checkError(err) {
			return nil, err
		}

		if res.Next() {
			err = res.Scan(&ownership.name)

			if !checkError(err) {
				return nil, err
			}
		}

		res, err = dataBase.Query("SELECT price FROM price_logs WHERE (timestamp = (SELECT MAX(timestamp) FROM price_logs))")

		if !checkError(err) {
			return nil, err
		}

		if res.Next() {
			err = res.Scan(&ownership.costPerOne)

			if !checkError(err) {
				return nil, err
			}
		}

		result = append(result, ownership)
	}

	return result, nil
}
