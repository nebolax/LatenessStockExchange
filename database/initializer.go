package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"strings"
)

var tableNameList = [7]string{"users", "stocks", "user_stock_ownerships",
	"price_logs", "transaction_logs", "comes_in", "event_logs"}

func checkError(err error) bool{
	if err != nil {
		panic(err)
		return false
	}
	return true
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

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
				tablesNames = remove(tablesNames, index)
				break
			}
		}
	}

	return len(tablesNames) == 0
}

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
