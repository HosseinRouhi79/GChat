package routers

import (
	"chat/config"
	"chat/src/api/handlers"
	"chat/src/api/middlewares/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.RouterGroup) {
	cfg := config.GetConfig()
	handler := handlers.AuthMobile{}

	handler2 := handlers.AuthRegister{}
	handler3 := handlers.UserHandler{}
	tokenHandler := handlers.TokenHandler{}
	r.POST("/register-login-mobile", handler.RLMobile)
	r.GET("/claim", middlewares.Authentication(cfg), tokenHandler.GetClaims)
	r.POST("/register", handler2.Register)
	r.POST("/send-otp", handler3.SendOtp)
	// r.POST("/refresh", handler.Refresh)
}
func Auth2(r *gin.RouterGroup) {
	handler0 := handlers.AuthUp{}
	r.POST("/up-login", handler0.UserPassLogin)
	// r.POST("/refresh", handler.Refresh)
}
