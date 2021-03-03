package netcommunicator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/nebolax/LatenessStockExcahnge/database/models"

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

type loginStatus string

const (
	//LoginOK is OK status
	LoginOK loginStatus = "Registration successfull"

	//UserUnexists fails when user with such nickname isn't registered
	UserUnexists loginStatus = "User with such login doesn't exists"

	//IncorrectPassword is thrown when password doesn't match the pwd in database
	IncorrectPassword loginStatus = "Incorrect password"
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
		tmpl, _ := template.ParseFiles("./templates/all-stocks-observer.html")
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

	//err := database.AddUser(login, email, pwd)

	// if !general.CheckError(err) {
	// 	processRegisterError(err)
	// 	return userRegFail
	// }

	// return newUserConfirmed

	if _, ok := users[login]; ok {
		return userRegFail
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		checkerr(err)
		users[login] = string(hash)
		return newUserConfirmed
	}
}

func loginUser(login, email, inpPwd string) (models.User, loginStatus) {
	if originPwd, ok := users[login]; !ok {
		return models.User{}, UserUnexists
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(originPwd), []byte(inpPwd))
		if err != nil {
			return models.User{}, IncorrectPassword
		} else {
			return models.User{}, LoginOK
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	logOutUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func portfolio(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		tmpl, _ := template.ParseFiles("./templates/portfolio.html")
		tmpl.Execute(w, "")
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
	user, status := loginUser(inpLogin, inpEmail, inpPwd)
	user = user
	tmpl, _ := template.ParseFiles("./templates/login.html")
	switch status {
	case LoginOK:
		setUserInfo(w, r, 1, inpLogin)
		http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
	case UserUnexists:
		tmpl.Execute(w, "user does not exist")
	case IncorrectPassword:
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
