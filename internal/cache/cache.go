package cache

import (
	"encoding/json"
	"shortify/internal/models"

	"github.com/imJayanth/go-modules/cache"
	"github.com/imJayanth/go-modules/config"
	"github.com/imJayanth/go-modules/errors"
)

func RedisGetUrl(appConfig *config.AppConfig, url *models.Url) errors.RestAPIError {
	keys, err := cache.RedisGetKeys(appConfig, url.GetRedisKey())
	if err.IsNotNull() {
		return err
	}
	for _, key := range keys {
		urlStr, err := cache.RedisGet(appConfig, key)
		if err.IsNotNull() {
			return err
		}
		json.Unmarshal([]byte(urlStr), &url)
		if url.Id > 0 {
			break
		}
	}
	if url.Id == 0 {
		return errors.NewNotFoundError("Record not found in redis")
	}
	return errors.NO_ERROR()
}

func RedisSetUrl(appConfig *config.AppConfig, url *models.Url) errors.RestAPIError {
	if err := cache.RedisSet(appConfig, url.GetRedisKey(), int(url.ExpiryTimeInHrs())*60*60, url.ToJson()); err.IsNotNull() {
		return err
	}
	return errors.NO_ERROR()
}
