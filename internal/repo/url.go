package repo

import (
	"shortify/internal/models"

	"github.com/imJayanth/go-modules/config"
	"github.com/imJayanth/go-modules/errors"
	"github.com/imJayanth/go-modules/helpers"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type UrlRepo struct {
	DB      *gorm.DB
	appConf *config.AppConfig
	name    string
}

func NewUrlRepo(appConfig *config.AppConfig) *UrlRepo {
	return &UrlRepo{
		DB:      appConfig.DBConfig.DB,
		appConf: appConfig,
		name:    "UrlRepo",
	}
}

func (ur *UrlRepo) SaveUrl(url *models.Url) errors.RestAPIError {
	m := "SaveUrl"
	if err := ur.DB.Table(url.TableName()).Save(url).Error; err != nil {
		ur.logError(m, "Save Error", []zapcore.Field{zap.String(url.TableName(), err.Error())})
		return errors.NewInternalServerError("Unable to save URL")
	}
	ur.logInfo(m, "Complete", []zapcore.Field{zap.Any("url", url)})
	return errors.NO_ERROR()
}

func (ur *UrlRepo) GetValidUrlByShortUrl(url *models.Url) errors.RestAPIError {
	m := "GetValidUrlByShortUrl"
	if err := ur.DB.Table(url.TableName()).Where(&models.Url{ShortUrl: url.ShortUrl}).Where("expires_at > ?", url.ExpiresAt).Find(url).Error; err != nil {
		ur.logError(m, "Get Error", []zapcore.Field{zap.String(url.TableName(), err.Error())})
		return errors.NewInternalServerError("Unable to get URL")
	}
	ur.logInfo(m, "Complete", []zapcore.Field{zap.Any("url", url)})
	return errors.NO_ERROR()
}

func (ur *UrlRepo) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(ur.appConf.LoggerConfig.Logger, ur.name, method, message, keyValues)
}

func (ur *UrlRepo) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(ur.appConf.LoggerConfig.Logger, ur.name, method, message, keyValues)
}
