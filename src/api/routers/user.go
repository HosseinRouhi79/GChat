package routers

import (
	"chat/src/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetOtp(r *gin.RouterGroup) {
	handler := handlers.UserHandler{}
	r.POST("/send-otp", handler.SendOtp)
}
func GetAllUsers(r *gin.RouterGroup) {
	r.GET("/users", handlers.GetAllUsers)
}


