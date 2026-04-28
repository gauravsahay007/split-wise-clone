package handler

import (
	"net/http"

	"github.com/gauravsahay007/split-wise-clone/auth"
	"github.com/gin-gonic/gin"
)

func (h *Handler) OAuthHandler(c *gin.Context) {
	provider := auth.OAuthProvider(c.Param("provider"))
	config := auth.GetOAuthConfig(provider)
	if config == nil {
		c.JSON(400, gin.H{"error": "Unsupported provider"})
		return
	}

	//You generate a random string (state)
	// Send it in the auth request
	// Store it (usually in session / cookie)
	// When Google/Github redirects back → it sends the same state
	// You verify it matches
	url := config.AuthCodeURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GenerateTokenFromGoogle(c *gin.Context) {
	code := c.Query("code")
	provider := auth.OAuthProvider(c.Param("provider"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	tokenObj, err := h.Service.OAuthCallback(code, provider)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenObj)
}
