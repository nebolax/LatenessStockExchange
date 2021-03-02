package database

// Main database error object
type DatabaseError struct {
	Message string
}

func (dbErr DatabaseError) Error() string {
	return dbErr.Message
}
