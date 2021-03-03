package netcommunicator

import (
	"fmt"
	"github.com/nebolax/LatenessStockExcahnge/database/models"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/nebolax/LatenessStockExcahnge/database"
	"github.com/nebolax/LatenessStockExcahnge/general"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/nebolax/LatenessStockExcahnge/pricesmonitor"
)

//OutcomingMessage is a struct
type OutcomingMessage struct {
	Type        string  `json:"type"`
	OffersCount int     `json:"offersCount"`
	StockPrice  float64 `json:"stockPrice"`
}

//incomingMessage is struct too
type incomingMessage struct {
	OfferType string `json:"offerType"`
}
type graphPageSetup struct {
	Type      string    `json:"type"`
	History   []float64 `json:"history"`
	CurOffers int       `json:"offers"`
}

type oneStockInfo struct {
	ID       int
	Name     string
	CurPrice float64
}

type regStatus int

const (
	newUserConfirmed regStatus = 0
	userRegFail      regStatus = 1
)

type loginStatus int

const (
	loginConfirmed        loginStatus = 0
	userUnexistsFail      loginStatus = 1
	incorrectPasswordFail loginStatus = 2
)

var (
	users     = make(map[string]string)
	clients   = make(map[*websocket.Conn]int)
	broadcast = make(chan int)
	key       = []byte("super-secret-key")
	store     = sessions.NewCookieStore(key)

	calcsNames = map[int]string{
		0: "Denis",
		1: "Serzh",
		2: "Leva",
		3: "Pasha",
		4: "Ilya",
		5: "ShchMax",
		6: "YurNik",
		7: "Nastya F",
		8: "Nastya Sh",
		9: "Nastya Ch",
	}
)

func getUserInfo(r *http.Request) *models.User {
	session, _ := store.Get(r, "user-info")
	id, ok := session.Values["userid"].(int)

	if !ok || id == 0 {
		return nil
	}

	result, err := database.GetUser(id)

	if !general.CheckError(err) {
		return nil
	}

	return result
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

func setUserInfo(w http.ResponseWriter, r *http.Request, id int, nickname string) {
	session, _ := store.Get(r, "user-info")
	session.Values["userid"] = id
	session.Save(r, w)
}

func logOutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-info")
	session.Values["userid"] = 0
	session.Save(r, w)
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	clients[ws] = id
	defer ws.Close()

	if calc := pricesmonitor.CalcByID(id); calc != nil {
		fmt.Printf("Id: %d, Offs: %d, st: %#v\n", calc.ID, calc.CurHandler.CurOffers, calc.History)
		ws.WriteJSON(graphPageSetup{"setup", calc.History, calc.CurHandler.CurOffers})
	}

	//Reading data from clients
	for {
		var msg incomingMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		var inc int
		if msg.OfferType == "sell" {
			inc = -1
		} else if msg.OfferType == "buy" {
			inc = 1
		}
		for _, calc := range pricesmonitor.AllCalculators {
			if calc.ID == id {
				calc.IncOffers(inc)
			}
		}
		fmt.Printf("New message: %s\n", msg.OfferType)
	}
}

//UpdateData is func
func UpdateData(id int, message OutcomingMessage) {
	for client := range clients {
		if clients[client] == id {
			err := client.WriteJSON(message)
			if err != nil {
				defer client.Close()
				delete(clients, client)
			}
		}
	}
}

func allStocksObserver(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		tmpl, _ := template.ParseFiles("./templates/all-stocks-observer.xhtml")
		var ar []oneStockInfo
		for _, calc := range pricesmonitor.AllCalculators {
			ar = append(ar, oneStockInfo{ID: calc.ID, Name: calcsNames[calc.ID], CurPrice: math.Round(calc.CurHandler.CurStock*100) / 100})
		}
		tmpl.Execute(w, ar)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func graphStockPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	id++
	file, _ := ioutil.ReadFile("./templates/graph-page.html")
	w.Write(file)
}

func processRegisterError(err error) {
	fmt.Println(err)
}

func regUser(login, email, pwd string) regStatus {

	err := database.AddUser(login, email, pwd)

	if !general.CheckError(err) {
		processRegisterError(err)
		return userRegFail
	}

	return newUserConfirmed

	/*if _, ok := users[login]; ok {
		return userRegFail
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		checkerr(err)
		users[login] = string(hash)
		return newUserConfirmed
	}*/
}

func loginUser(login, email, inpPwd string) loginStatus {
	if originPwd, ok := users[login]; !ok {
		return userUnexistsFail
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(originPwd), []byte(inpPwd))
		if err != nil {
			return incorrectPasswordFail
		} else {
			return loginConfirmed
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	logOutUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func showError(w http.ResponseWriter, r *http.Request, err NetError){
	tmpl, _ := template.ParseFiles("./templates/error.html")
	tmpl.Execute(w, err)
}

func portfolio(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		userInfo := getUserInfo(r)

		if userInfo == nil {
			showError(w, r, NetError{"User not found!"})
			return
		}

		tmpl, _ := template.ParseFiles("./templates/portfolio.html")
		tmpl.Execute(w, userInfo)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func procRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inpLogin := r.PostForm.Get("login")
	inpPwd := r.PostForm.Get("password")
	inpEmail := r.PostForm.Get("email")
	status := regUser(inpLogin, inpEmail, inpPwd)
	switch status {
	case newUserConfirmed:
		setUserInfo(w, r, 1, inpLogin)
		http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
	case userRegFail:
		tmpl, _ := template.ParseFiles("./templates/register.html")
		tmpl.Execute(w, "User already exists")
	}
}

func getRegisterHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/register.html")
	tmpl.Execute(w, "")
}

func procLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inpLogin := r.PostForm.Get("login")
	inpPwd := r.PostForm.Get("password")
	inpEmail := r.PostForm.Get("email")
	status := loginUser(inpLogin, inpEmail, inpPwd)
	tmpl, _ := template.ParseFiles("./templates/login.html")
	switch status {
	case loginConfirmed:
		setUserInfo(w, r, 1, inpLogin)
		http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
	case userUnexistsFail:
		tmpl.Execute(w, "user does not exist")
	case incorrectPasswordFail:
		tmpl.Execute(w, "incorrect password")
	}
}

func getLoginHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/login.html")
	tmpl.Execute(w, "")
}

//StartServer is func
func StartServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./templates"))))

	router.HandleFunc("/graph{id:[0-9]+}", graphStockPage)
	router.HandleFunc("/allstocks", allStocksObserver)
	router.HandleFunc("/portfolio", portfolio)
	router.HandleFunc("/login", procLogin).Methods("POST")
	router.HandleFunc("/login", getLoginHTML).Methods("GET")
	router.HandleFunc("/register", procRegister).Methods("POST")
	router.HandleFunc("/register", getRegisterHTML).Methods("GET")
	router.HandleFunc("/logout", logout).Methods("POST")

	router.HandleFunc("/ws/graph{id:[0-9]+}", handleConnections)
	http.Handle("/", router)

	log.Println("starting http server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
