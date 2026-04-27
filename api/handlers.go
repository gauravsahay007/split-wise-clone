package api

import (
	"net/http"
	"strconv"

	"github.com/gauravsahay007/split-wise-clone/business"
	"github.com/gauravsahay007/split-wise-clone/models"
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

// @Summary Create a new expense
// @Description Add an expense to a group (requires authentication)
// @Tags Expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Expense true "Expense Data"
// @Success 201 {object} map[string]string "Expense added"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /expenses [post]
func (h *Handler) ExpenseHandler(c *gin.Context) {
	val, exists := c.Get("current_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context missing"})
		return
	}
	currentUserID := val.(int)
	var exp models.Expense
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exp.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}
	if len(exp.UserIDs) == 0 && len(exp.Shares) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No participants provided"})
		return
	}

	exp.PaidBy = currentUserID
	if err := h.Service.CreateExpense(exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Expense added"})
}

// @Summary Create a new group
// @Description Create a group with current user as owner
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Group Request" example({"name":"Trip Group"})
// @Success 201 {object} models.Group
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Failed to create group"
// @Router /groups [post]
func (h *Handler) CreateGroupHandler(c *gin.Context) {
	userID := c.MustGet("current_user_id").(int)
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	group, err := h.Service.CreateGroup(req.Name, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create group"})
		return
	}
	c.JSON(http.StatusCreated, group)
}

// @Summary Get group balances
// @Description Get simplified balances for a group
// @Tags Balances
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {array} models.Balance
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /groups/{id}/balances [get]
func (h *Handler) BalancesHandler(c *gin.Context) {
	// Get group_id from URL: /api/groups/:id/balances
	groupIDStr := c.Param("id")
	groupID, _ := strconv.Atoi(groupIDStr)

	balances, err := h.Service.GetBalances(groupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, balances)
}

// @Summary Add member to group
// @Description Add a user to a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body object true "User ID" example({"user_id":2})
// @Success 200 {object} map[string]string "User added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Failed to add member"
// @Router /groups/{id}/members [post]
func (h *Handler) AddMemberHandler(c *gin.Context) {
	// Get group_id from URL /api/groups/:id/members
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req struct {
		UserID int `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := h.Service.AddMemberToGroup(groupID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add user to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to group successfully"})
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
