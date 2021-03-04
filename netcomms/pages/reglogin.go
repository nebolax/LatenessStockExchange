package pages

import (
	"html/template"
	"net/http"

	"github.com/nebolax/LatenessStockExcahnge/general/status"
	"github.com/nebolax/LatenessStockExcahnge/netcomms/procs"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	procs.LogoutUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ProcRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inpLogin := r.PostForm.Get("login")
	inpPwd := r.PostForm.Get("password")
	inpEmail := r.PostForm.Get("email")
	id, cs := procs.RegUser(inpLogin, inpEmail, inpPwd)
	switch cs {
	case status.OK:
		procs.SetSessionUserID(w, r, id)
		http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
	case status.UserRegFail:
		tmpl, _ := template.ParseFiles("./templates/register.html")
		tmpl.Execute(w, "User already exists")
	}
}

func GetRegisterHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/register.html")
	tmpl.Execute(w, "")
}

func ProcLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inpLogin := r.PostForm.Get("login")
	inpPwd := r.PostForm.Get("password")
	//inpEmail := r.PostForm.Get("email")
	id, cs := procs.LoginUser(inpLogin, inpPwd)
	tmpl, _ := template.ParseFiles("./templates/login.html")
	switch cs {
	case status.OK:
		procs.SetSessionUserID(w, r, id)
		http.Redirect(w, r, "/portfolio", http.StatusSeeOther)
	case status.NoSuchUser:
		tmpl.Execute(w, "user does not exist")
	case status.IncorrectPassword:
		tmpl.Execute(w, "incorrect password")
	}
}

func SendLoginHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/login.html")
	tmpl.Execute(w, "")
}
