package controllers

import (
	"Medods/internal/database"
	"Medods/internal/email"
	"Medods/models"
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

	if len(guid) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GUID provided, the length must be <= 128"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate hash" + err.Error()})
		return
	}

	client := models.User{Email: body.Email, Password: string(hash), GUID: guid}

	accessToken, _, err := JWTMaker.CreateToken(client.ID, client.Email, c.ClientIP(), 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create an access token: " + err.Error()})
		return
	}

	refreshToken, err := models.NewRefreshToken()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a refresh token: " + err.Error()})
		return
	}

	hashToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate hash " + err.Error()})
		return
	}
	client.RefreshToken = string(hashToken)

	result := database.DB.Create(&client)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func Login(c *gin.Context) {
	var body struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body: " + err.Error()})
		return
	}

	var client models.User
	result := database.DB.Where("email = ?", body.Email).First(&client)

	CheckIP()

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

	if err := bcrypt.CompareHashAndPassword([]byte(client.RefreshToken), []byte(body.RefreshToken)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token: " + err.Error()})
		return
	}

	RenewTokens(c, body.RefreshToken, client.Email)
}

func CheckIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is empty"})
			return
		}
		claims, err := JWTMaker.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token"})
			return
		}
		currentIP := c.ClientIP()

		if claims.IP != currentIP {
			go func() {
				err := email.SendEmailWarning(claims.Email, claims.IP, currentIP)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "error occurred while sending email"})
					return
				}
			}()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your IP address has changed. You have to reauthenticate"})
			return
		}
		c.Next()
	}
}

func RenewTokens(c *gin.Context, refreshToken, email string) {

	var user models.User

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(refreshToken)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, _, err := JWTMaker.CreateToken(user.ID, user.Email, c.ClientIP(), 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an access token: " + err.Error()})
		return
	}

	newRefreshToken, err := models.NewRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a refresh token: " + err.Error()})
		return
	}

	hashToken, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate hash " + err.Error()})
		return
	}

	database.DB.Model(&user).Updates(models.User{
		RefreshToken: string(hashToken),
	})

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken, "refreshToken": newRefreshToken})
}
