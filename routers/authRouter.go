package router

import ("github.com/gin-gonic/gin"
		"go-jwt-model/controller"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/users/signup",controller.Signup())
	incomingRoutes.POST("users/login",controller.Login())
}