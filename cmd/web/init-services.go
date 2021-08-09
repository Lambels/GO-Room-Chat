package main

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/handlers"
	"GO-Group-Chat/internal/helpers"
	"GO-Group-Chat/internal/middleware"
)

func initServices(app *config.AppConfig) {
	handlers.InitHandlers(app)
	helpers.InitHelpers(app)
	middleware.InitializeMiddleware(app)
}