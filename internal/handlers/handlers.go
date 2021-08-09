package handlers

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/helpers"
	"log"
	"net/http"
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
		log.Println(err)
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

		err := helpers.Signup(w, r, email, usn, []byte(password))
		switch err {
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