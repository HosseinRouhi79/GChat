package handlers

import (
	"chat/config"
	"chat/data/cache"
	"chat/data/db"
	"chat/data/models"
	"chat/services"
	"chat/src/api/helper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OtpRequest struct {
	Mobile string `form:"mobile" binding:"mobile,min=11,max=11"`
}

type UserHandler struct {
	UserService services.UserService
	OtpService  services.OtpService
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all users from the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} helper.HTTPResponse "List of users"
// @Failure 500 {object} helper.HTTPResponse "Internal server error"
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	db := db.GetDB()

	// Define a slice to hold the list of users
	var users []models.User

	// Query the database to fetch all users
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Return the list of users as a JSON response
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetAllUsers godoc
// @Summary Get active users
// @Description Retrieves a list of all users from the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} helper.HTTPResponse "List of active users"
// @Failure 500 {object} helper.HTTPResponse "Internal server error"
// @Router /users/active [get]
func GetActiveUsers(c *gin.Context) {
	usernameList := []string{}
	redis := cache.GetRedis()
	// ctx := context.Background()

	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = redis.Scan(cursor, "*", 10).Result()
		if err != nil {
			fmt.Println("Error during SCAN:", err)
			return
		}

		for _, key := range keys {
			val, err := redis.Get(key).Result()
			if err != nil {
				fmt.Println("Error getting value for key", key, ":", err)
				continue
			}
			usernameList = append(usernameList, val)
		}

		c.JSON(http.StatusOK, gin.H{"active": usernameList})

		if cursor == 0 {
			break
		}
	}
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Deletes a user with the specified ID from the system.
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 404 {object} helper.HTTPResponse "User not found"
// @Failure 500 {object} helper.HTTPResponse "Internal server error"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	db := db.GetDB()
	id := c.Param("id")

	// Find the user by ID
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete related records in user_roles table
	if err := db.Where("user_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related user roles"})
		return
	}

	// Delete the user

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// OTP godoc
// @Summary Send OTP request
// @Description OTP request
// @Tags auth
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param mobile formData string true "mobile"
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 400 {object} helper.HTTPResponse "Failed"
// @Router /send-otp/ [post]
func (h UserHandler) SendOtp(c *gin.Context) {
	cfg := config.GetConfig()
	// logger := logging.NewLogger(cfg)
	h = UserHandler{
		OtpService:  *services.NewOtpService(cfg),
		UserService: *services.NewUserService(cfg),
	}
	req := OtpRequest{}
	err := c.ShouldBind(&req)
	var status bool = true
	var code int = 200
	if err != nil {
		h.UserService.Logger.Infof("error has occured: %s", err.Error())
		status = false
		code = 500
		c.JSON(code, gin.H{
			"status": status,
			"err":    err.Error(),
		})
		return
	}
	otpCode := helper.GenerateOtp()
	err = h.OtpService.SetOtp(req.Mobile, strconv.Itoa(otpCode))

	if err != nil {
		h.UserService.Logger.Infof("can not set otp code: %s", err.Error())
		status = false
		code = 400
		c.JSON(code, gin.H{
			"status": status,
			"err":    err.Error(),
		})
		return
	}

	// send OTP SMS
	c.JSON(code, gin.H{
		"status": status,
		"OTP":    strconv.Itoa(otpCode),
	})
}
