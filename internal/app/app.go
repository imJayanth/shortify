package app

import (
	"log"
	"strconv"

	"github.com/imJayanth/go-modules/config"
	"github.com/imJayanth/go-modules/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func StartApplication(appConfig *config.AppConfig) {
	var app = fiber.New()
	// Add logger middleware
	app.Use(func(c *fiber.Ctx) error {
		helpers.LogRequest(appConfig.LoggerConfig.Logger, c.Method(), c.Path())
		err := c.Next()
		helpers.LogResponse(appConfig.LoggerConfig.Logger, c.Response().StatusCode())
		return err
	})

	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(recover.New())
	mapUrls(app, appConfig)
	log.Fatal(app.Listen(":" + strconv.Itoa(appConfig.ServerConfig.APIPort)))
}
