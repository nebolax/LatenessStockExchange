package procs

import (
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/general/status"
)

//LogoutUser is func
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	SetSessionUserID(w, r, 0)
}

//RegUser is func
func RegUser(login, email, pwd string) (int, status.StatusCode) {

	id, err := database.AddUser(login, email, pwd)

	if !general.CheckError(err) {
		return 0, status.UserRegFail
	}

	return id, status.OK
}

//LoginUser is func
func LoginUser(login, inpPwd string) (int, status.StatusCode) {
	id, cs := database.LoginByNickname(login, inpPwd)

	if cs != status.OK {
		return 0, cs
	}
	return id, status.OK
}
