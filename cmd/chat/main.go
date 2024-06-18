package main

import (
	"chat-system/internal/api/cache"
	"chat-system/internal/api/routes"
	dbmanager "chat-system/internal/db_manager"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	csSession := dbmanager.InitDB()
	defer csSession.Close()

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
