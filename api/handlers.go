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

func (h *Handler) UserHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	user, err := h.Service.CreateUser(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) ExpenseHandler(c *gin.Context) {
	var exp models.Expense

	if err := c.ShouldBindBodyWithJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
