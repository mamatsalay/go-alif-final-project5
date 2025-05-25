package main

import (
	"log"
	"workout-tracker/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Error loading .env file")
		return
	}

	app.StartServer()
}
