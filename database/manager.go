// Package for working with database
package database

import "database/sql"

// Reference to the batadase object
var dataBase *sql.DB

// If the batadase was connected to the program this is true
var initialized = false
