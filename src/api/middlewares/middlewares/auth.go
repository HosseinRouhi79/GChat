package middlewares

import (
	"chat/config"
	"chat/services"
	"chat/src/constants"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var claimsMap map[string]interface{}
		tokenService := services.NewTokenService(cfg)
		auth := c.GetHeader(constants.Authorization)
		if auth == "" {
			err = fmt.Errorf("auth header doesnt exist")
		} else {
			token := strings.Split(auth, " ")
			accessToken := token[1]

			claimsMap, err = tokenService.GetClaims(accessToken)
			if err != nil {
				err = fmt.Errorf("error has occured:%s", err.Error())
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.Set("id", claimsMap["id"])
		c.Set("userName", claimsMap["userName"])
		c.Set("firstName", claimsMap["firstName"])
		c.Set("lastName", claimsMap["lastName"])
		c.Set("mobileNumber", claimsMap["mobileNumber"])
		c.Set("email", claimsMap["email"])
		c.Set("roles", claimsMap["roles"])
		c.Set("exp", claimsMap["exp"])

		c.Next()

	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var roleList []string
		roleval, ok := c.Keys["roles"]
		if !ok {
			err = fmt.Errorf("error has occured")
		} else {
			roleList = append(roleList, roleval.([]string)...)
		}
		if err == nil {
			for _, role := range roleList {
				for _, vRole := range validRoles {
					if role == vRole {
						c.Next()
					}
				}
			}
		}
	}
}
