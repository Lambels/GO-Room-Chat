package driver

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	SQL 	*sqlx.DB
}

// Connects to the mysql databse given dsnString
// then pings the databse to verify a live connection
//
// returns an sqlx database
func ConnectMySQL(dsn string) (*DB, error) {
	db, err := sqlx.Connect("mysql", dsn)	// Connects and check with ping

	return &DB{db}, err
}