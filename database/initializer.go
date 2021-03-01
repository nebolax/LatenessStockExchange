package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"io/ioutil"
	"strings"
)

// List of names of database tables. It is used to check if database is incomplete
var tableNameList = [7]string{"users", "stocks", "user_stock_ownerships",
	"price_logs", "transaction_logs", "comes_in", "event_logs"}

// Check tables for incompleteness. Returns true if everything is OK.
func checkTables(db *sql.DB) bool {
	tables, err := db.Query("SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';")
	if !checkError(err) {
		return false
	}
	tablesNames := tableNameList[:]
	for tables.Next(){
		var name string
		err := tables.Scan(&name)
		if !checkError(err) {
			return false
		}
		for index, value := range  tablesNames {
			if value == name {
				tablesNames = general.Remove(tablesNames, index)
				break
			}
		}
	}

	return len(tablesNames) == 0
}

// Generate tables from start (if smth is wrong)
func createTables(db *sql.DB) {
	fmt.Println("Creating new tables")
	file, err := ioutil.ReadFile("database\\storage\\template.sql")
	checkError(err)

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		result, err := db.Exec(request)
		checkError(err)
		fmt.Println(result)
	}
}

// Initialization of database. If something is wrong, recreate all database
func Init() {
	db, err := sql.Open("sqlite3", "database\\storage\\database.db")
	checkError(err)
	var dbOk = checkTables(db)
	if !dbOk {
		createTables(db)
	}

	initialized = true
	dataBase = db
}
