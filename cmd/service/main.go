package main

import (
	"Medods/internal/database"
	"Medods/routes"
)

func init() {
	database.ConnectToDB()
}

func main() {
	r := routes.NewRoutes()
	r.Run(":8080")
}
