package database

import "github.com/nebolax/LatenessStockExcahnge/general"

// Main database error object
type DatabaseError struct {
	message string
}

func (dbErr DatabaseError) Error() string {
	return dbErr.message
}

// Check error for null-ness (if there any exception)
func checkError(err error) bool{
	if err != nil {
		if general.DebugMode {
			panic(err)
		}
		return false
	}
	return true
}
