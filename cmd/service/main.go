package main

import (
	"Medods/internal/controllers"
	"Medods/internal/database"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateResfreshToker() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

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
	r.GET("/getTokens/:guid", func(c *gin.Context) {
		guid := c.Param("guid")
		c.JSON(http.StatusOK, gin.H{
			"guid": guid,
		})
	})
	r.POST("/signup", controllers.SignUp)
	r.GET("/login", controllers.Login)

	r.Run(":8080")
}
