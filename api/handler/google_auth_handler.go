package handler

import (
	"net/http"

	"github.com/gauravsahay007/split-wise-clone/auth"
	"github.com/gin-gonic/gin"
)

// GoogleLoginHandler godoc
// @Summary Login with Google
// @Description Redirects user to Google OAuth consent screen
// @Tags Auth
// @Success 307 {string} string "Redirect to Google"
// @Router /auth/google [get]
func (h *Handler) GoogleLoginHandler(c *gin.Context) {
	url := auth.GoogleConfig().AuthCodeURL("random-state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GenerateTokenFromGoogle godoc
// @Summary Google OAuth Callback
// @Description Handles Google callback, creates or logs in user, returns JWT token
// @Tags Auth
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} map[string]interface{} "JWT Token"
// @Failure 400 {object} map[string]string "Missing code"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/google/callback [get]
func (h *Handler) GenerateTokenFromGoogle(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	tokenObj, err := h.Service.GoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenObj)
}
