package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"

	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/driver"
	"GO-Group-Chat/internal/models"

	"github.com/srinathgs/mysqlstore"
)

func SetupApp() (*config.AppConfig) {
	dbName := flag.String("dbName", "testdb", "Database name")
	dbUser := flag.String("dbUser", "root", "Databse user")
	dbPass := flag.String("dbPass", "", "Database password")
	dbPort := flag.String("dbPort", "3306", "Database port")
	dbHost := flag.String("dbHost", "localhost", "Database host")
	stName := flag.String("stName", "sessions", "Sessions table name")
	port := flag.String("port", ":8080", "Port on wich the server listens on")
	
	flag.Parse()

	log.Println("setupApp: Initializing database")
	var dsnString string

	if *dbPass == "" {
		dsnString = fmt.Sprintf("%s:@tcp(%s:%s)/%s?parseTime=true&loc=Local", *dbUser, *dbHost, *dbPort, *dbName)
	} else {
		dsnString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", *dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	}

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
		ActiveConnections: make(map[int64]models.CommunicationChannels),
	}
}