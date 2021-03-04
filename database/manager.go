// Package for working with database
package database

import (
	"database/sql"
	"sync"
)

type Database struct {
	db *sql.DB
	mu sync.Mutex
}

func (obj *Database) Query (request string) (*sql.Rows, error) {
	obj.mu.Lock()
	res, err := obj.db.Query(request)
	obj.mu.Unlock()

	return res, err
}

func (obj *Database) Exec (request string) (sql.Result, error) {
	obj.mu.Lock()
	res, err := obj.db.Exec(request)
	obj.mu.Unlock()

	return res, err
}

// Reference to the batadase object
var dataBase Database

// If the batadase was connected to the program this is true
var initialized = false
