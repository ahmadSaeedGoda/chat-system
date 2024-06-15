package main

import (
	"chat-system/internal/api/cache"
	"chat-system/internal/api/routes"
	"chat-system/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	db.Init()
	defer db.Session.Close()

	cache.Init()

	r := routes.InitRoutes()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), r))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}
}
