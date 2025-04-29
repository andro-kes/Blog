package controllers

import (
	"github.com/andro-kes/Blog/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

var oauth2Config *oauth2.Config

func init() {
	config.LoadConfig()
	oauth2Config = &oauth2.Config{
		ClientID: config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		RedirectURL: config.REDIRECT_URL,
		Scopes: []string{"login:email"},
		Endpoint: yandex.Endpoint,
	}
}