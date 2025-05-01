package controllers

import (
	"Medods/internal/database"
	"Medods/models/user"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body: " + err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate hash" + err.Error()})
		return
	}
	client := user.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&client)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": client})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body: " + err.Error()})
		return
	}

	var client user.User
	result := database.DB.Where("email = ?", body.Email).First(&client)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find user: " + result.Error.Error()})
		return
	}
	if client.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email address: " + client.Email})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password: " + err.Error()})
		return
	}
	token, err := CreateJWT(client.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create token: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

func CreateJWT(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
