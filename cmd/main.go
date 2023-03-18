package main

import (
	"shortify/internal/app"
	"shortify/internal/models"

	"github.com/imJayanth/go-modules/config"
)

func main() {
	appConfig := config.SetupConfig()
	appConfig.DBConfig.DBName = "shortify_db"
	appConfig.DBConfig.Automigrate = map[string]interface{}{
		models.TABLE_URLS: &models.Url{},
	}
	appConfig.SetupDatabase()
	appConfig.SetupLogger()
	appConfig.SetupRedis()
	app.StartApplication(appConfig)
}
