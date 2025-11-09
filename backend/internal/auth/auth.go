package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	maxAge = 86400 * 30
	isProd = false
)

func NewAuth() {
	googleClientId := os.Getenv("G_CLIENT_ID")
	googleClientSecret := os.Getenv("G_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(os.Getenv("G_KEY")))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	var redirectURL string
	if isProd {
		redirectURL = "https://kiitwallet.dev"
	} else {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}
	gothic.Store = store
	goth.UseProviders(google.New(googleClientId, googleClientSecret, redirectURL))

}
