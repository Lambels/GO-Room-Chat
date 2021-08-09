package main

import (
	"GO-Group-Chat/internal/routes"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	Version = "1.0.0"
)

func main() {
	app := SetupApp()
	initServices(app)

	defer func(){
		app.DB.SQL.Close()
		app.Store.Close()
	}()

	log.Println("--------------------------------------------")
	log.Printf("Running GO-Group-Chat on v%s", app.Version)
	log.Printf("Running with %d processors", runtime.NumCPU())
	log.Printf("Running on %s", runtime.GOOS)
	log.Println("--------------------------------------------")

	srv := http.Server {
		Addr: app.Domain,
		Handler: routes.InitHandlers(),
		IdleTimeout:       	30 * time.Second,
		ReadTimeout:       	10 * time.Second,
		ReadHeaderTimeout: 	5 * time.Second,
		WriteTimeout:      	5 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", app.Domain)
	log.Fatal(srv.ListenAndServe())
}