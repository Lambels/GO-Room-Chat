package helpers

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/models"
	"database/sql"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	app *config.AppConfig

	ErrUserNotFound	= errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrAlreadyLoggedIn = errors.New("user already logged in")
	ErrNotLoggedIn = errors.New("user not logged in")
)

// InitHelpers populates the global app variable making it available to the root of the package
func InitHelpers(a *config.AppConfig) {
	app = a
}

// IsAuth checks if the user is authentificated and return true if he is and false if he isnt
func IsAuth(r *http.Request) (bool) {
	session, err := app.Store.Get(r, "session")
	if err != nil {
		return false
	}

	a, ok := session.Values["active"]
	if !ok {
		return false
	}
	
	return a.(bool)
}

// Login logs the user in and checks everything
//
// Checks:
//
// Account exists,
// User isnt logged in,
// Passwords match
func Login(w http.ResponseWriter, r *http.Request, email string, pass []byte) (error) {
	var account models.Account

	// Already logged in?
	if IsAuth(r) {
		return ErrAlreadyLoggedIn
	}

	switch err := app.DB.SQL.Get(&account, "SELECT * FROM accounts WHERE Email = ?", email); err {
	case nil:

	case sql.ErrNoRows:
		return ErrUserNotFound

	default:
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), pass); err != nil {
		return err
	}

	session, err := app.Store.New(r, "session")
	if err != nil {
		return err
	}

	session.Values["active"] = true
	session.Values["APK"] = account.ID
	return session.Save(r, w)
}

// Logout logs the user out and checks everything
//
// Checks:
//
// User is logged in
func Logout(w http.ResponseWriter, r *http.Request) (error) {
	// Not logged in?
	if !IsAuth(r) {
		return ErrNotLoggedIn
	}

	session, err := app.Store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["active"] = false

	return session.Save(r, w)
}

// Signup logs a user in and checks everything
//
// Checks:
//
// Account is unique,
// User isnt already logged in
func Signup(w http.ResponseWriter, r *http.Request, email, usn string, pass []byte) (error) {
	var account models.Account
	if IsAuth(r) {
		return ErrAlreadyLoggedIn
	}

	switch err := app.DB.SQL.Get(&account, "SELECT * FROM accounts WHERE email = ?", email); err {
	case nil:
		return ErrUserAlreadyExists

	case sql.ErrNoRows:

	default:
		return err
	}

	hashedPass, _ := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)

	res := app.DB.SQL.MustExec("INSERT INTO accounts (email, username, password) VALUES (?, ?, ?)", email, usn, hashedPass)
	apk, _ := res.LastInsertId()

	session, err := app.Store.New(r, "session")
	if err != nil {
		return err
	}

	session.Values["active"] = true
	session.Values["APK"] = apk
	return session.Save(r, w)
} 

// Helper function for template.ExecuteTemplate: https://pkg.go.dev/html/template#Template.ExecuteTemplate
func RenderTemplate(w io.Writer, name string, data interface{}) (error) {
	return app.Templates.ExecuteTemplate(w, name, data)
}

// GetUser returns the user based that it exists and return ErrUserNotFound if there is no user
func GetUser(r *http.Request) (models.Account, error) {
	var acc models.Account

	session, err := app.Store.Get(r, "session")
	if err != nil {
		return acc, err
	}

	if _, ok := session.Values["active"]; !ok || !session.Values["active"].(bool) {
		return acc, ErrUserNotFound
	}

	apk := session.Values["APK"]

	if err := app.DB.SQL.Get(&acc, "SELECT * FROM accounts WHERE ID = ?", apk); err != nil {
		return acc, err
	}

	return acc, nil
}