package handlers

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/helpers"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	app *config.AppConfig
)

func InitHandlers(a *config.AppConfig) {
	app = a
}

// Handles "/" and renders to different templates depending if the user
// sending the request is auth
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

func Logout(w http.ResponseWriter, r *http.Request) {
	err := helpers.Logout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}