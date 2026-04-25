package api

import (
	"net/http"
	"strconv"

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
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	user, err := h.Service.CreateUser(req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ExpenseHandler adds a new expense to a specific group
func (h *Handler) ExpenseHandler(c *gin.Context) {
	var exp models.Expense

	// ShouldBindJSON maps the "group_id", "paid_by", etc., from JSON to the struct
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Basic Validation: Ensure group and amount are provided
	if exp.GroupID == 0 || exp.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id and a positive amount are required"})
		return
	}

	if err := h.Service.CreateExpense(exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense added successfully"})
}

func (h *Handler) CreateGroupHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	group, _ := h.Service.Repo.SaveGroup(req.Name)
	c.JSON(200, group)
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
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authorization header",
			})

			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer" {
			tokenString = tokenString[7:]
		}

		userID, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

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
