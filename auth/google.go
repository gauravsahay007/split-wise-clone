package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var clientId, clientSecret string

func LoadGoogleAuthEnv(client_id string, client_secret string) {
	clientId = client_id
	clientSecret = client_secret
}

func GoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
