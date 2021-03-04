package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nebolax/LatenessStockExcahnge/general"
)

const StandartPath = "database\\storage\\"
const testPath = "storage\\"

// List of names of database tables. It is used to check if database is incomplete
var tableNameList = [7]string{"users", "stocks", "user_stock_ownerships",
	"price_logs", "transaction_logs", "comes_in", "event_logs"}

// Check tables for incompleteness. Returns true if everything is OK.
func checkTables(db *sql.DB) bool {
	tables, err := db.Query("SELECT Name FROM sqlite_master WHERE type ='table' AND Name NOT LIKE 'sqlite_%';")
	if !general.CheckError(err) {
		print("CRINGE!!!\n" + err.Error())
		return false
	}

	defer tables.Close()

	tablesNames := make([]string, len(tableNameList))
	copy(tablesNames, tableNameList[:])

	for tables.Next() {
		var name string
		err := tables.Scan(&name)
		if !general.CheckError(err) {
			return false
		}
		for index, value := range tablesNames {
			if value == name {
				tablesNames = general.Remove(tablesNames, index)
				break
			}
		}
	}

	return len(tablesNames) == 0
}

// Generate tables from start (if smth is wrong)
func createTables(db *sql.DB, path string) {
	fmt.Println("Creating new tables")
	file, err := ioutil.ReadFile(path + "template.sql")
	general.CheckError(err)

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.Exec(request)
		general.CheckError(err)
		//fmt.Println(result)
	}
}

// Initialization of database. If something is wrong, recreate all database
func Init(path string) {
	db, err := sql.Open("sqlite3", path+"database.db")
	general.CheckError(err)
	var dbOk = checkTables(db)
	if !dbOk {
		createTables(db, path)
	}

	initialized = true
	dataBase = Database{db: db}
}
