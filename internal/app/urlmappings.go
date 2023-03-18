package app

import (
	"shortify/internal/routers"

	"github.com/imJayanth/go-modules/config"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App *fiber.App
}

func NewRouter(app *fiber.App) *Router {
	return &Router{App: app}
}

func mapUrls(app *fiber.App, appConfig *config.AppConfig) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	router := routers.NewRouter(app)
	router.SetUrlRoutes(appConfig)
}
