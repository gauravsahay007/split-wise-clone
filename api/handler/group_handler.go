package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) GetUserGroupsHandler(c *gin.Context) {
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

	groups, err := h.Service.FetchUserGroups(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, groups)
}

func (h *Handler) GetGroupMembers(c *gin.Context) {
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

	gid, err := strconv.Atoi(c.Query("groupId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid groupId"})
		return
	}

	res, err := h.Service.GetGroupMembers(userID, gid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}
