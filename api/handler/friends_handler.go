package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gauravsahay007/split-wise-clone/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleAddFriend(c *gin.Context) {
	val, exists := c.Get("current_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User context missing"})
		return
	}

	friendIdsString := c.Query("friends")
	if friendIdsString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No friend found to be added"})
	}

	friendsIds := utils.ParseCSVToString(friendIdsString)

	userID, ok := val.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error: Invalid user ID format"})
		return
	}

	err := h.Service.AddFriend(userID, friendsIds)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": "Friend Added Successfully"})
}

func (h *Handler) GetFriendsList(c *gin.Context) {
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

	res, err := h.Service.GetFriendsList(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, res)
}

func (h *Handler) SearchFriendsInAGroup(c *gin.Context) {
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

	searchString := c.Query("search")
	gid, err := strconv.Atoi(c.Query("groupId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid groupId"})
		return
	}
	fmt.Println(searchString)
	res, err := h.Service.SearchFriendsInAGroup(userID, searchString, gid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, res)
}
