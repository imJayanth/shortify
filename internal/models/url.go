package models

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/imJayanth/go-modules/helpers"
)

const (
	TABLE_URLS             = "urls"
	URL_EXPIRY_TIME_IN_HRS = 24
)

func (Url) TableName() string {
	return TABLE_URLS
}

type Url struct {
	Id          int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	OriginalUrl string    `json:"original_url" gorm:"column:original_url" validate:"required,url"`
	Domain      string    `json:"domain" gorm:"column:domain"`
	ShortUrl    string    `json:"short_url" gorm:"column:short_url"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:DATETIME"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:DATETIME"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"column:expires_at;type:DATETIME"`
}

func (url *Url) Shorten() {
	// Generate a SHA256 hash of the long URL
	hash := sha256.Sum256([]byte(url.OriginalUrl))
	// Convert the hash to a base64-encoded string
	base64Hash := base64.URLEncoding.EncodeToString(hash[:])
	// Use the base 62 algorithm to convert the base64-encoded hash to a short URL
	url.ShortUrl = helpers.Base62Encode(base64Hash)
}

func (url *Url) ExpiryTimeInHrs() time.Duration {
	return time.Duration(URL_EXPIRY_TIME_IN_HRS)
}

func (url *Url) GetRedisKey() string {
	replacer := strings.NewReplacer("!@#$^&*()-=!@#$^&*()-=", "!@#$^&*()-=")
	return replacer.Replace(fmt.Sprintf(`!@#$^&*()-=%v!@#$^&*()-=%v!@#$^&*()-=`, url.OriginalUrl, url.ShortUrl))
}

func (url *Url) ToJson() string {
	js, serr := json.Marshal(url)
	if serr != nil {
		log.Println("Error while marshalling: ", serr)
	}
	return string(js)
}
