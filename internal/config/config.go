package config

import (
	"GO-Group-Chat/internal/driver"
	"html/template"

	"github.com/srinathgs/mysqlstore"
)	

type AppConfig struct {
	Templates 		*template.Template
	Store			*mysqlstore.MySQLStore
	DB 				*driver.DB
	Domain 			string
	Version			string
}