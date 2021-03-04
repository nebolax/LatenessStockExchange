package database

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/general/status"
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
	var template = regexp.MustCompile("[[:alpha:]]+@179.ru")
	return template.Match([]byte(email))
}

// Adds user with given email and password and nickname to batadase
func AddUser(nickname string, email string, password string) (int, error) {
	if !checkEmail(email) {
		return 0, DatabaseError{"Email address is incorrect!"}
	}

	salt := randomString(saltLength)
	resultPassword := password + hashKey + salt
	hasher := sha512.New()
	hasher.Write([]byte(resultPassword))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if !initialized {
		return 0, DatabaseError{"Database is not initialized yet!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO users (username, email, password_hash, password_salt, money) values ('%s', '%s', '%s', '%s', 0);",
		nickname, email, hash, salt))

	if !general.CheckError(err) {
		return 0, err
	}

	result, err := dataBase.Query(fmt.Sprintf("SELECT id FROM users WHERE (username = '%s')", nickname))

	if !general.CheckError(err) {
		return 0, err
	}

	defer result.Close()

	if result.Next() {
		var id int
		err = result.Scan(&id)

		return id, err
	}

	return 0, DatabaseError{"Some cringy cringe - there is no user with just created nickname"}
}

// Tries to login user from sql query using given password. Returns error if something went wrong
// Returns nil if everything is OK
func login(user *sql.Rows, password string) (int, status.StatusCode) {
	if user.Next() {
		var salt string
		var passwordHash string
		var id int

		err := user.Scan(&salt, &passwordHash, &id)

		if !general.CheckError(err) {
			return 0, status.DatabaseError
		}

		resultString := password + hashKey + salt

		hasher := sha512.New()
		hasher.Write([]byte(resultString))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		if hash != passwordHash {
			return 0, status.IncorrectPassword
		} else {
			return id, status.OK
		}

	} else {
		return 0, status.NoSuchUser
	}
}

// Tries to login user with given nickname by given password
// Returns error if something went wrong
// Returns nil if everything is OK
func LoginByNickname(nickname string, password string) (int, status.StatusCode) {
	if !initialized {
		return 0, status.DatabaseNotInited
	}

	query, err := dataBase.Query(fmt.Sprintf(
		"SELECT password_salt, password_hash, id FROM users WHERE (username = '%s')", nickname))

	if !general.CheckError(err) {
		return 0, status.DatabaseError
	}

	defer query.Close()

	return login(query, password)
}

// Tries to login user with given email by given password
// Returns error if something went wrong
// Returns nil if everything is OK
func LoginByEmail(email string, password string) (int, status.StatusCode) {
	if !initialized {
		return 0, status.DatabaseNotInited
	}

	query, err := dataBase.Query(fmt.Sprintf(
		"SELECT password_salt, password_hash FROM users WHERE (email = '%s')", email))

	if !general.CheckError(err) {
		return 0, status.DatabaseError
	}

	defer query.Close()

	return login(query, password)
}

// Returns map from StockId to Amount - ownerships of userId
// If error occurred, returns nil as map and error as second object
func GetInvestmentPortfolio(userId int) (map[int]int, error) {
	if !initialized {
		return nil, DatabaseError{"Database is not initialized"}
	}

	query, err := dataBase.Query(fmt.Sprintf(
		"SELECT stock_id, amount FROM user_stock_ownerships WHERE (user_id = %d)", userId))

	if !general.CheckError(err) {
		return nil, err
	}

	defer query.Close()

	var result = make(map[int]int)

	for query.Next() {
		var stockId int
		var amount int

		err = query.Scan(&stockId, &amount)

		if !general.CheckError(err) {
			return nil, err
		}

		result[stockId] = amount
	}

	return result, nil

}

// Returns more info about investment portfolio of userId (see Ownership for more information)
func GetInvestmentPortfolioPretty(userId int) ([]Ownership, error) {
	data, err := GetInvestmentPortfolio(userId)
	if !general.CheckError(err) {
		return nil, err
	}

	var result = make([]Ownership, 0)

	for key, value := range data {
		var ownership = Ownership{key, "", value, 0}

		res, err := dataBase.Query(fmt.Sprintf("SELECT Name FROM stocks WHERE (id = %d)", key))
		if !general.CheckError(err) {
			return nil, err
		}

		defer res.Close()

		if res.Next() {
			err = res.Scan(&ownership.Name)

			if !general.CheckError(err) {
				return nil, err
			}
		}

		res1, err := dataBase.Query("SELECT price FROM price_logs WHERE (timestamp = (SELECT MAX(timestamp) FROM price_logs))")

		if !general.CheckError(err) {
			return nil, err
		}

		defer res1.Close()

		if res1.Next() {
			err = res1.Scan(&ownership.CostPerOne)

			if !general.CheckError(err) {
				return nil, err
			}
		}

		result = append(result, ownership)
	}

	return result, nil
}

func GetUser(id int) (*User, error) {
	if !initialized {
		return nil, DatabaseError{"Database is not initialized"}
	}

	rows, err := dataBase.Query(fmt.Sprintf("SELECT username, email, money FROM users WHERE (id = '%d')", id))
	if !general.CheckError(err) {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var result User = User{}
		result.Id = id
		err = rows.Scan(&result.Nickname, &result.Email, &result.Money)

		if !general.CheckError(err) {
			return nil, err
		}

		result.Stocks, err = GetInvestmentPortfolioPretty(result.Id)

		if !general.CheckError(err) {
			return nil, err
		}

		return &result, nil
	}

	return nil, DatabaseError{"No such user!"}
}
