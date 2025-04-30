package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/beevik/guid"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("medods-company") // os.Getenv("SECRET_KEY")

func getRole(user string) string {
	if user == "admin" {
		return "admin"
	}
	return "user"
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
	fmt.Println("Start working!")
	g := guid.New()
	fmt.Println(g.String())
}
