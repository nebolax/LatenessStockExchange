package procs

import (
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
)

//RegStatus is status
type RegStatus int

const (
	//NewUserConfirmed is status
	NewUserConfirmed RegStatus = 0

	//UserRegFail is status too
	UserRegFail RegStatus = 1
)

//LoginStatus is status
type LoginStatus string

const (
	//LoginOK is OK status
	LoginOK LoginStatus = "Registration successfull"

	//UserUnexists fails when user with such nickname isn't registered
	UserUnexists LoginStatus = "User with such login doesn't exists"

	//IncorrectPassword is thrown when password doesn't match the pwd in database
	IncorrectPassword LoginStatus = "Incorrect password"
)

//LogoutUser is func
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	SetSessionUserID(w, r, 0)
}

//RegUser is func
func RegUser(login, email, pwd string) (int, RegStatus) {

	id, err := database.AddUser(login, email, pwd)

	if !general.CheckError(err) {
		return 0, UserRegFail
	}

	return id, NewUserConfirmed
}

//LoginUser is func
func LoginUser(login, inpPwd string) (int, LoginStatus) {
	id, err := database.LoginByNickname(login, inpPwd)

	if !general.CheckError(err) {
		return 0, LoginStatus(err.Error())
	}
	return id, LoginOK
}
