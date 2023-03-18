package services

import (
	"shortify/internal/cache"
	"shortify/internal/models"
	"shortify/internal/repo"

	"time"

	"github.com/imJayanth/go-modules/config"
	"github.com/imJayanth/go-modules/errors"
	"github.com/imJayanth/go-modules/helpers"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type UrlService struct {
	urlRepo *repo.UrlRepo
	appConf *config.AppConfig
	name    string
}

func NewUrlService(appConfig *config.AppConfig) *UrlService {
	return &UrlService{
		urlRepo: repo.NewUrlRepo(appConfig),
		appConf: appConfig,
		name:    "UrlService",
	}
}

func (us *UrlService) GenerateShortUrl(url *models.Url) errors.RestAPIError {
	m := "GenerateShortUrl"
	us.logInfo(m, "Start", nil)

	if err := cache.RedisGetUrl(us.appConf, url); err.IsNull() {
		us.logInfo(m, "Complete", []zapcore.Field{zap.String("redis value", url.ShortUrl)})
		return errors.NO_ERROR()
	}

	url.Shorten()
	url.CreatedAt = us.appConf.CurrentTime()
	url.ExpiresAt = url.CreatedAt.Add(url.ExpiryTimeInHrs() * time.Hour)
	if err := us.urlRepo.SaveUrl(url); err.IsNotNull() {
		return err
	}

	if err := cache.RedisSetUrl(us.appConf, url); err.IsNotNull() {
		us.logError(m, "Redis", []zapcore.Field{zap.String("error", err.Message)})
		return err
	}
	us.logInfo(m, "Complete", nil)
	return errors.NO_ERROR()
}

func (us *UrlService) GetValidUrlByShortUrl(url *models.Url) errors.RestAPIError {
	m := "GetValidUrlByShortUrl"
	us.logInfo(m, "Start", nil)
	defer us.logInfo(m, "Complete", nil)

	if err := cache.RedisGetUrl(us.appConf, url); err.IsNull() {
		return errors.NO_ERROR()
	}

	if err := us.urlRepo.GetValidUrlByShortUrl(url); err.IsNotNull() {
		return err
	}
	us.logInfo(m, "Complete", nil)
	return errors.NO_ERROR()
}

func (us *UrlService) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(us.appConf.LoggerConfig.Logger, us.name, method, message, keyValues)
}

func (us *UrlService) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(us.appConf.LoggerConfig.Logger, us.name, method, message, keyValues)
}
