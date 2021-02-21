package database

type DatabaseError struct {
	message string
}

func (dbErr DatabaseError) Error() string {
	return dbErr.message
}
