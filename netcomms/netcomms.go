package netcommunicator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nebolax/LatenessStockExcahnge/netcomms/pages"
)

func landingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reached")
	http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/graph{id:[0-9]+}", pages.GraphStockPage)
	router.HandleFunc("/allstocks", pages.AllStocksObserver)
	router.HandleFunc("/portfolio", pages.Portfolio)
	router.HandleFunc("/login", pages.ProcLogin).Methods("POST")
	router.HandleFunc("/login", pages.SendLoginHTML).Methods("GET")
	router.HandleFunc("/register", pages.ProcRegister).Methods("POST")
	router.HandleFunc("/register", pages.GetRegisterHTML).Methods("GET")
	router.HandleFunc("/logout", pages.Logout).Methods("POST")
	router.HandleFunc("/", landingPage)
}

//StartServer is func
func StartServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./templates"))))

	setupRoutes(router)

	router.HandleFunc("/ws/graph{id:[0-9]+}", pages.HandleConnections)
	http.Handle("/", router)

	log.Println("starting http server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
