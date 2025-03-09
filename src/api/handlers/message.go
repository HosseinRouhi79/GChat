package handlers

import (
	"chat/config"
	"chat/data/db"
	"chat/data/models"
	"chat/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeleteMessage(c *gin.Context) {
	db := db.GetDB()
	id := c.Param("id")

	// Find the message by ID
	var message models.Message
	if err := db.First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Delete the message
	var cfg = config.GetConfig()
	tokenService := services.NewTokenService(cfg)
	token := c.GetHeader("Authorization")
	accessToken := strings.Split(token, " ")
	tk := accessToken[1]
	mpClaims, _ := tokenService.GetClaims(tk)
	fmt.Println(mpClaims)
	if int(mpClaims["id"].(float64)) != message.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	} else if err := db.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
