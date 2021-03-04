package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs"
)

//Portfolio func
func Portfolio(w http.ResponseWriter, r *http.Request) {
	if procs.IsUserLoggedIn(r) {
		userInfo, err := procs.GetSessionUserInfo(r)
		if !general.CheckError(err) {
			fmt.Println(err)
			return
		}

		tmpl, err := template.ParseFiles("./templates/portfolio.html")
		checkerr(err)
		tmpl.Execute(w, userInfo)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
