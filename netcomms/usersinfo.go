package netcommunicator

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/general/status"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func getUserInfo(r *http.Request) (*database.User, status.StatusCode) {
	session, _ := store.Get(r, "user-info")
	id, ok := session.Values["userid"].(int)

	if !ok || id == 0 {
		return nil, status.NotLoggedIn
	}

	result, err := database.GetUser(id)

	if !general.CheckError(err) {
		return nil, status.DatabaseError
	}

	return result, status.OK
}

func isUserLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, "user-info")
	id, ok := session.Values["userid"].(int)

	if !ok || (id == 0) {
		return false
	} else {
		return true
	}
}

func setUserInfo(w http.ResponseWriter, r *http.Request, id int) {
	session, _ := store.Get(r, "user-info")
	session.Values["userid"] = id
	session.Save(r, w)
}
