package handler

import (
	"net/http"

	"github.com/gauravsahay007/split-wise-clone/business"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *business.Service
}

// @Summary Register a new user
// @Description Create a new user with name, password, email, and profile picture
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object true "User Request" example({"name":"Gaurav","password":"123456","email":"gaurav@example.com","profile_pic":"https://img.com/pic.png"})
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /users [post]
func (h *Handler) UserHandler(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		Password   string `json:"password" binding:"required,min=6"`
		Email      string `json:"email"`
		ProfilePic string `json:"profile_pic"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.CreateUser(req.Name, req.Password, req.Email, req.ProfilePic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object true "Login Request" example({"id":1,"password":"123456"})
// @Success 200 {object} map[string]string "JWT Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req struct {
		ID       int    `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "ID and Passwor required",
		})
		return
	}

	token, err := h.Service.Authenticate(req.ID, req.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Unauthorized: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

// @Summary Get user summary
// @Description Get overall balance summary of the logged-in user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/summary [get]
func (h *Handler) UserSummaryHandler(c *gin.Context) {
	val, exists := c.Get("current_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User context missing"})
		return
	}
	userID, ok := val.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error: Invalid user ID format"})
		return
	}

	summary, err := h.Service.GetUserOverallSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate user summary: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
