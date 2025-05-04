package routes

import (
	"Medods/internal/controllers"

	"github.com/gin-gonic/gin"
)

func NewRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Next()

		if c.Writer.Header().Get("Content-Type") == "" {
			c.Header("Content-Type", "application/json")
		}
	})

	api := r.Group("/api")
	{
		api.GET("/refresh", controllers.Login)
		api.POST("/getTokens/:guid", controllers.SignUp)
	}
	return r
}
