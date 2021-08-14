package routes

import (
	"net/http"
	"GO-Group-Chat/internal/handlers"
	"GO-Group-Chat/internal/middleware"

	multiplexer "github.com/gorilla/mux"
)

func InitHandlers() http.Handler {
	mux := multiplexer.NewRouter()

	authProtected := mux.NewRoute().Subrouter()
	authProtected.Use(middleware.OnlyAuthMiddleware)
	unAuthProtected := mux.NewRoute().Subrouter()
	unAuthProtected.Use(middleware.OnlyUnAuthMiddleware)

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	mux.HandleFunc("/", handlers.Index)

	authProtected.HandleFunc("/index/ws", handlers.IndexWS)

	authProtected.HandleFunc("/logout", handlers.Logout).Methods(http.MethodGet)
	
	authProtected.HandleFunc("/rooms/new", handlers.CreateRoom).Methods(http.MethodGet, http.MethodPost)
	authProtected.HandleFunc("/rooms/{pk:[0-9]+}", handlers.JoinRoom).Methods(http.MethodGet)
	authProtected.HandleFunc("/rooms/{pk:[0-9]+}/ws", handlers.JoinRoomWS)
	// TODO: Implement rooms and handlers
	authProtected.HandleFunc("/rooms/delete/{pk}", nil).Methods(http.MethodGet, http.MethodDelete)

	unAuthProtected.HandleFunc("/signup", handlers.Signup).Methods(http.MethodGet, http.MethodPost)
	unAuthProtected.HandleFunc("/login", handlers.Login).Methods(http.MethodGet, http.MethodPost)

	return mux
}