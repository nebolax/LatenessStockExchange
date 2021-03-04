package procs

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//GetSessionUserInfo is func
func GetSessionUserInfo(r *http.Request) (*database.User, error) {
	id := GetSessionUserID(r)

	result, err := database.GetUser(id)

	if !general.CheckError(err) {
		return nil, err
	}

	return result, nil
}

//IsUserLoggedIn is func
func IsUserLoggedIn(r *http.Request) bool {
	id := GetSessionUserID(r)
	if id == 0 {
		return false
	}
	return true
}

//SetSessionUserID id func to set user id:D!!
func SetSessionUserID(w http.ResponseWriter, r *http.Request, id int) {
	session, _ := store.Get(r, "user-info")
	session.Values["userid"] = id
	session.Save(r, w)
}

//GetSessionUserID is func
func GetSessionUserID(r *http.Request) int {
	session, _ := store.Get(r, "user-info")
	id, ok := session.Values["userid"].(int)
	if !ok {
		return 0
	}
	return id
}
