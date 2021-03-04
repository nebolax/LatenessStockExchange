package pages

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nebolax/LatenessStockExcahnge/general"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
	} else {
		ws, err := websocket.Upgrade(w, r, nil, 0, 0)
		general.CheckError(err)

		vars := mux.Vars(r)
		graphID, _ := strconv.Atoi(vars["id"])

		if calc := pricesmonitor.CalcByID(graphID); calc != nil {
			procs.RegNewGraphObserver(ws, calc, procs.GetSessionUserID(r))
		}
	}
}

func GraphStockPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/graph-page.html")
	general.CheckError(err)
	if procs.IsUserLoggedIn(r) {
		tmpl.Execute(w, "loggedIn")
	} else {
		tmpl.Execute(w, "guest")
	}
}
