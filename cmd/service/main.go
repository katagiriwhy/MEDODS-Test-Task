package main

import (
	"Medods/models/employee"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("medods-company") // os.Getenv("SECRET_KEY")

func getRole(user string) string {
	if user == "admin" {
		return "admin"
	}
	return "employee"
}

func CreateResfreshToker() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func CreateJWT(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": username,
		"aud": getRole(username),
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func main() {
	r := gin.Default()
	r.GET("/getTokens/:guid", func(c *gin.Context) {
		guid := c.Param("guid")
		c.JSON(http.StatusOK, gin.H{
			"guid": guid,
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	employee.CreateEmployee(r)
	r.Run(":8080")
}
