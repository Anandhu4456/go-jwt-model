package routers

import (
	"go-jwt/controller"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)


func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/user/:user_id",controller.GetUser())
}