package status

//StatusCode handles all types of errors which can happen
type StatusCode string

const (
	OK StatusCode = "Operation completed successfully"

	UserRegFail StatusCode = "Couldn't register user"

	NoSuchUser StatusCode = "User with such login doesn't exists"

	IncorrectPassword StatusCode = "Incorrect password"

	DatabaseNotInited StatusCode = "Database isn't initialized"

	NotLoggedIn StatusCode = "User isn't logged in"

	DatabaseError StatusCode = "Unknown error happened in the database"
)
