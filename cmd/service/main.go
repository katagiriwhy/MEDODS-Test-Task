package main

import (
	"Medods/internal/database"
	"Medods/routes"
	"log"
)

func init() {
	database.ConnectToDB()
}

func main() {
	r := routes.NewRoutes()
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
