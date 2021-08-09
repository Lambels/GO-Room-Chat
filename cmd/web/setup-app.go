package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"

	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/driver"

	"github.com/srinathgs/mysqlstore"
)

func SetupApp() (*config.AppConfig) {
	dbName := flag.String("dbName", "testdb", "Database name")
	stName := flag.String("stName", "sessions", "Sessions table name")
	port := flag.String("port", ":8080", "Port on wich the server listens on")
	
	flag.Parse()

	dsnString := fmt.Sprintf("root:@tcp(localhost:3306)/%s?parseTime=true&loc=Local", *dbName)

	log.Println("setupApp: Initializing database")
	db, err := driver.ConnectMySQL(dsnString)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("setupApp: Initializing sessions")
	store, err := mysqlstore.NewMySQLStoreFromConnection(db.SQL.DB, *stName, "/", 3600, []byte("super-secret-key"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("setupApp: Parsing templates")
	return &config.AppConfig {
		Templates: template.Must(template.ParseGlob("./templates/*.html")),
		Store: store,
		DB: db,
		Domain: *port,
		Version: Version,
	}
}