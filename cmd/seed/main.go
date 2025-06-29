package main

import (
	"auth_service/database/seeders"
	"auth_service/internal/config"
	"fmt"
	"log"
)

func main() {
	config.LoadENV()
	db := config.ConnectDB()

	if err := seeder.SeedRun(db); err != nil {
		log.Fatalf("seeding gagal: %v", err)
	}

	fmt.Println("seeder berhasil...")
}
