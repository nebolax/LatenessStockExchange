package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func clearTest() {
	_, _ = dataBase.Exec("DROP TABLE *")
}

func TestCheckTablesOk(t *testing.T) {
	db, err := sql.Open("sqlite3", "storage\\database.db")

	if err != nil {
		t.Error(err.Error())
		return
	}

	if !checkTables(db) {
		t.Error("CheckTables is not ok")
	}

	_ = db.Close()
}

func TestCheckTablesBad(t *testing.T) {
	_, _ = os.Create("test.sqlite")
	defer os.Remove("test.sqlite")

	db, err := sql.Open("sqlite3", "test.sqlite")

	if err != nil {
		t.Error(err.Error())
	}

	if checkTables(db) {
		t.Error("CheckTables is not ok")
	}
}

func TestCreateTables(t *testing.T) {
	_, _ = os.Create("test.sqlite")
	defer os.Remove("test.sqlite")

	db, err := sql.Open("sqlite3", "test.sqlite")
	if err != nil {
		t.Error(err.Error())
	}

	createTables(db)

	if !checkTables(db) {
		t.Error("DB is not ok after creating")
	}
}

func TestInit(t *testing.T) {
	Init(testPath)

	if !initialized {
		t.Error("DB is not initialized after Init()")
	}

	if !checkTables(dataBase) {
		t.Error("DB is not OK after Init()")
	}
}
