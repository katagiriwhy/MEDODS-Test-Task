package controllers

import (
	"Medods/internal/database"
	"Medods/models"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var JWTMaker = models.NewTokenMaker(os.Getenv("JWT_SECRET"))

func SignUp(c *gin.Context) {

	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body: " + err.Error()})
		return
	}
	guid := c.Param("guid")
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No GUID provided"})
		return
	}
	if len(guid) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GUID provided"})
		return
	}
	fmt.Println(body.Email, body.Password)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate hash" + err.Error()})
		return
	}
	client := models.User{Email: body.Email, Password: string(hash), GUID: guid}
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
	var client models.User
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
	token, claims, err := JWTMaker.CreateToken(client.ID, client.Email, c.ClientIP(), 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create token: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

func CreateJWT(userID uint, userIP string) (string, error) {

}

func CreateResfreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
