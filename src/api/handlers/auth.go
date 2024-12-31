package handlers

import (
	"chat/config"
	"chat/dto"
	"chat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var cfg = config.GetConfig()

type TokenHandler struct {
	Token string `form:"token"`
}
type AuthMobile struct {
	Mobile string `form:"mobile" binding:"mobile,min=11,max=11,required"`
	Otp    string `form:"otp"`
}
type AuthRegister struct {
	FirstName string `form:"firstName" binding:"required, min=3"`
	LastName  string `form:"lastName" binding:"required, min=3"`
	Username  string `form:"username" binding:"required, min=3"`
	Email     string `form:"email" binding:"required, email"`
	Password  string `form:"password" binding:"required, password"`
}

// User_Auth godoc
// @Summary User Auth
// @Description Register
// @Tags auth
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param firstName formData string true "first name"
// @Param lastName formData string true "last name"
// @Param username formData string true "username"
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 400 {object} helper.HTTPResponse "Failed"
// @Router /register/ [post]
func (re AuthRegister) Register(c *gin.Context) {

	err := c.ShouldBind(&re)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := dto.RegisterationDto{
		FirstName: re.FirstName,
		LastName:  re.LastName,
		Username:  re.Username,
		Email:     re.Email,
		Password:  re.Password,
	}

	userService := services.NewUserService(cfg)
	err = userService.Register(dto)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Registration Done",
	})

}

// User_Auth godoc
// @Summary User Auth
// @Description Register Login
// @Tags auth
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param mobile formData string true "Mobile number"
// @Param otp formData string false "OTP"
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 400 {object} helper.HTTPResponse "Failed"
// @Router /register-login-mobile/ [post]
func (auth AuthMobile) RLMobile(c *gin.Context) {

	am := AuthMobile{}
	err := c.ShouldBind(&am)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := dto.RegisterLoginByMobileDto{
		Mobile: am.Mobile,
		Otp:    am.Otp,
	}

	userService := services.NewUserService(cfg)
	tokenDetail, err := userService.RegisterLoginByMobile(dto)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"token":  tokenDetail,
	})
}

// User_Auth godoc
// @Summary Get Claims
// @Description Register Login
// @Tags auth
// @Security AuthBearer
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param token formData string false "jwt token"
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 400 {object} helper.HTTPResponse "Failed"
// @Router /claim/ [post]
func (t TokenHandler) GetClaims(c *gin.Context) {
	err := c.ShouldBind(&t)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokenService := services.NewTokenService(cfg)
	// fmt.Println(t.Token)
	mpClaims, _ := tokenService.GetClaims(t.Token)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"claims": mpClaims,
	})
}
