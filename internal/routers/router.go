package routers

import "github.com/gofiber/fiber/v2"

type Router struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) *Router {
	return &Router{app: app}
}
