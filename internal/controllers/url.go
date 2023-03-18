package controllers

import (
	"fmt"
	"shortify/internal/models"
	"shortify/internal/services"
	"strings"

	"github.com/imJayanth/go-modules/config"
	"github.com/imJayanth/go-modules/errors"
	"github.com/imJayanth/go-modules/helpers"

	"github.com/imJayanth/go-modules/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/go-playground/validator.v9"
)

type UrlController struct {
	urlService *services.UrlService
	appConf    *config.AppConfig
	name       string
}

func NewUrlController(appConfig *config.AppConfig) *UrlController {
	return &UrlController{
		urlService: services.NewUrlService(appConfig),
		appConf:    appConfig,
		name:       "UrlController",
	}
}

func (uc *UrlController) GenerateShortUrl(c *fiber.Ctx) error {
	m := "GenerateShortUrl"
	uc.logInfo(m, "Start", nil)
	url := models.Url{}
	if err := c.BodyParser(&url); err != nil {
		return err
	}
	validate := validator.New()
	if e := validate.Struct(url); e != nil {
		return response.RespondBadRequest(c, e.Error())
	}
	domain := c.Hostname()
	if !strings.Contains(domain, "localhost:") {
		if port := c.Port(); port != "" {
			domain += ":" + port
		}
	}
	url.Domain = domain
	if strings.Contains(url.OriginalUrl, url.Domain) {
		return response.RespondBadRequest(c, "cannot shorten a shortened URL")
	}
	url.ShortUrl = ""
	uc.logInfo(m, "Input", []zapcore.Field{zap.String("original url", url.OriginalUrl)})
	if saveErr := uc.urlService.GenerateShortUrl(&url); saveErr.IsNotNull() {
		return response.RespondError(c, saveErr)
	}
	uc.logInfo(m, "Complete", nil)
	return response.RespondCreated(c, map[string]string{
		"url": fmt.Sprintf("%v/%v", url.Domain, url.ShortUrl),
	})
}

func (uc *UrlController) LookupShortUrl(c *fiber.Ctx) error {
	m := "LookupShortUrl"
	uc.logInfo(m, "Start", nil)
	url := models.Url{}
	url.ShortUrl = c.Params("SHORTURL")
	if strings.TrimSpace(url.ShortUrl) == "" {
		err := errors.NewBadRequestError("Invalid SHORTURL")
		uc.logError(m, "Params", []zapcore.Field{zap.String("SHORTURL", err.Message)})
		return response.RespondError(c, err)
	}

	uc.logInfo(m, "Input", []zapcore.Field{zap.String("short url", url.ShortUrl)})
	if getErr := uc.urlService.GetValidUrlByShortUrl(&url); getErr.IsNotNull() {
		return response.RespondError(c, getErr)
	}
	uc.logInfo(m, "Complete", nil)
	c.Set("Location", url.OriginalUrl)
	return c.Status(fiber.StatusMovedPermanently).SendString("")
}

func (uc *UrlController) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(uc.appConf.LoggerConfig.Logger, uc.name, method, message, keyValues)
}

func (uc *UrlController) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(uc.appConf.LoggerConfig.Logger, uc.name, method, message, keyValues)
}
