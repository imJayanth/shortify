package routers

import (
	"shortify/internal/controllers"

	"github.com/imJayanth/go-modules/config"
)

func (r *Router) SetUrlRoutes(appConfig *config.AppConfig) {
	urlController := controllers.NewUrlController(appConfig)
	api := r.app.Group("/")
	api.Post("/", urlController.GenerateShortUrl)
	api.Get("/:SHORTURL", urlController.LookupShortUrl)
}
