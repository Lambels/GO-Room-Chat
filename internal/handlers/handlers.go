package handlers

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/helpers"
	"GO-Group-Chat/internal/middleware"
	"GO-Group-Chat/internal/models"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var (
	app *config.AppConfig

	upgrader = websocket.Upgrader {
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}

	regChan 	chan<- models.WsConnection
	unRegChan 	chan<- models.WsConnection
	broadChan 	chan<- *models.Message
)

func init() {
	regChan, unRegChan, broadChan = models.NewConnStore()
}

func InitHandlers(a *config.AppConfig) {
	app = a
}

// Handles "/" and renders to different templates depending if the user
// sending the request is auth
//
// and accepts only the GET verb
func Index(w http.ResponseWriter, r *http.Request) {
	if !helpers.IsAuth(r) {
		err := helpers.RenderTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Couldn't render template", http.StatusInternalServerError)
		}
		return
	}

	u, err := helpers.GetUser(r)
	if err != nil {
		return
	}

	err = helpers.RenderTemplate(w, "index.html", u)
	if err != nil {
		http.Error(w, "Couldn't render template", http.StatusInternalServerError)
	}
}

// Handels "/index/ws" and connects to the index room used to feed users on the index page updates such as
// new rooms
//
// IndexWS needs to be protected by onlyAuth middleware
func IndexWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	// browser doesent support ws ?
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	user := r.Context().Value(middleware.AccountKey{})

	conn := models.WsConnection {
		Conn: 	ws,
		Acc: 	user.(models.Account),
		Send: 	make(chan *models.Message),
	}

	go conn.Listen(regChan, unRegChan, broadChan)
}

// Signup handles "/signup" and accepts POST and GET verbs
//
// signup needs to be protected by onlyUnAuth middleware
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.FormValue("email")
		usn := r.FormValue("username")
		password := r.FormValue("password")

		switch err := helpers.Signup(w, r, email, usn, []byte(password)); err {
		case nil:
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		
		case helpers.ErrUserAlreadyExists:
			http.Redirect(w, r, "/auth-err/user-exists", http.StatusSeeOther)
			return
		
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	helpers.RenderTemplate(w, "signup.html", nil)
}

// Login handles "/login" and accepts POST and GET verbs
//
// login needs to be protected by onlyUnAuth middleware
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.FormValue("email")
		password := r.FormValue("password")

		switch err := helpers.Login(w, r, email, []byte(password)); err {
		case nil:
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		
		case helpers.ErrUserNotFound:
			http.Redirect(w, r, "/auth-err/user-not-found", http.StatusSeeOther)
			return

		case bcrypt.ErrMismatchedHashAndPassword:
			http.Redirect(w, r, "/auth-err/wrong-pass", http.StatusSeeOther)
			return

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	helpers.RenderTemplate(w, "login.html", nil)
}

// Logout handles "/logout" and accepts the GET verb
//
// Logout needs to be protected by onlyAUth middleware
func Logout(w http.ResponseWriter, r *http.Request) {
	err := helpers.Logout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}