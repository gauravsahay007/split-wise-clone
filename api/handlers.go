package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gauravsahay007/split-wise-clone/business"
	"github.com/gauravsahay007/split-wise-clone/models"
	"github.com/gauravsahay007/split-wise-clone/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *business.Service
}

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

// ExpenseHandler adds a new expense to a specific group
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

// BalancesHandler calculates the net settlements with simplification
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Missing or invalid token format",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 4. Set the ID in context for other handlers to use
		c.Set("current_user_id", userID)

		c.Next()
	}
}

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
