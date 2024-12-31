package handlers

import (
	"chat/config"
	"chat/services"
	"chat/src/api/helper"
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
