package models

import "github.com/nebolax/LatenessStockExcahnge/database"

type User struct {
	Id int
	Nickname string
	Email string
	Money float64
	Stocks []database.Ownership
}
