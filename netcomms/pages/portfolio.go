package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/general/status"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs"
)

//Portfolio func
func Portfolio(w http.ResponseWriter, r *http.Request) {
	if procs.IsUserLoggedIn(r) {
		userInfo, cs := procs.GetSessionUserInfo(r)
		if cs != status.OK {
			fmt.Println(cs)
			return
		}

		tmpl, err := template.ParseFiles("./templates/portfolio.html")
		general.CheckError(err)
		tmpl.Execute(w, userInfo)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
