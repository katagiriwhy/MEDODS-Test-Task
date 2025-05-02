package main

import (
	"Medods/internal/controllers"
	"Medods/internal/database"

	"github.com/gin-gonic/gin"
)

func init() {
	database.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Next()

		if c.Writer.Header().Get("Content-Type") == "" {
			c.Header("Content-Type", "application/json")
		}
	})
	r.POST("/getTokens/:guid", controllers.SignUp)
	r.GET("/refresh", controllers.Login)

	r.Run(":8080")
}
