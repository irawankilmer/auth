package main

import (
	"auth_service/database"
	"auth_service/internal/config"
	"log"
)

func main() {
	config.LoadENV()
	db := config.ConnectDB()

	if err := database.RunMigration(db); err != nil {
		log.Fatalf("Migrasi gagal: %v", err)
	}
}
