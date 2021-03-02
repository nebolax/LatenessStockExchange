package models

import "github.com/nebolax/LatenessStockExcahnge/database"

type User struct {
	id int
	nickname string
	email string
	money float64
	stocks []database.Ownership
}
