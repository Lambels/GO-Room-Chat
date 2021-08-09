package middleware

import (
	"context"
	"net/http"
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/helpers"
	"GO-Group-Chat/internal/models"
)

var app *config.AppConfig

type AccountKey struct{}

func InitializeMiddleware(a *config.AppConfig) {
	app = a 
}

// OnlyAuthMiddleware redirects to "/register" un auth users.
//
// If auth user is found the account is stored under the r.Context()
// with the key AccountKey{}
func OnlyAuthMiddleware(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect un auth users to "/register"
		if !helpers.IsAuth(r) {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		// Get Account
		var account *models.Account
		session, _ := app.Store.Get(r, "session")
		apk := session.Values["APK"]
		app.DB.SQL.Get(account, "SELECT * FROM accounts WHERE ID = ?", apk.(string))
		
		// Store account under r.Context()
		ctx := context.WithValue(r.Context(), AccountKey{}, account)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// OnlyUnAuthMiddleware redirect auth users to "/"
func OnlyUnAuthMiddleware(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect to "/" auth users
		if helpers.IsAuth(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}