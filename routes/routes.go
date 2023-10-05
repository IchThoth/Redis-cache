package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	routes.GET("/users/profile")
}
