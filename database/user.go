package database

type User struct {
	Id       int
	Nickname string
	Email    string
	Money    float64
	Stocks   []Ownership
}
