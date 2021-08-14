package config

import (
	"GO-Group-Chat/internal/driver"
	"GO-Group-Chat/internal/models"
	"html/template"

	"github.com/srinathgs/mysqlstore"
)	

type AppConfig struct {
	Templates 			*template.Template
	Store				*mysqlstore.MySQLStore
	DB 					*driver.DB
	Domain 				string
	Version				string
	ActiveConnections	map[int64]models.CommunicationChannels
}