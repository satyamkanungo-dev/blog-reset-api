package main

import (
	"fmt"
	"log"

	"github.com/satyamkanungo-dev/blog-rest-api/internal/config"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// database
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	defer pool.Close()

	fmt.Println("Everything is ok")
}
