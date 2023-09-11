package main

import (

	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	r := gin.New()
	r.Use(gin.Logger())

	

	r.GET("/user", func(c *gin.Context) {
		c.JSON(200,gin.H{"message":"hey"})
		
	})
	r.Run(":" + port)
}
