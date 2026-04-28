package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthProvider string

type OAuthUser struct {
	Provider   string
	ProviderID string
	Email      string
	Name       string
	Picture    string
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type GithubUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

const (
	Google   OAuthProvider = "google"
	Github   OAuthProvider = "github"
	LinkedIn OAuthProvider = "linkedin"
)

var googleClientId, googleClientSecret, githubClientId, githubClientSecret string

func LoadGoogleAuthEnv(client_id string, client_secret string) {
	googleClientId = client_id
	googleClientSecret = client_secret
}

func LoadGithubAuthEnv(client_id string, client_secret string) {
	githubClientId = client_id
	githubClientSecret = client_secret
}

func GetOAuthConfig(provider OAuthProvider) *oauth2.Config {
	switch provider {
	case Google:
		return &oauth2.Config{
			ClientID:     googleClientId,
			ClientSecret: googleClientSecret,
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}

	case Github:
		return &oauth2.Config{
			ClientID:     githubClientId,
			ClientSecret: githubClientSecret,
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
	}

	return nil
}
