package routers

import (
	"chat/src/api/handlers"

	"github.com/gin-gonic/gin"
)

func RouteMessage(r *gin.RouterGroup){

	r.DELETE("/:id/delete", handlers.DeleteMessage)
}