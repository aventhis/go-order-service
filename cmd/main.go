package main

import (
	"fmt"
	"github.com/aventhis/go-order-service/internal/config"
	"github.com/aventhis/go-order-service/internal/db"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database := db.Init(cfg.Database)
	defer database.Close()

	fmt.Println("App running on port:", cfg.Server.Port)
	fmt.Println("Connected to DB:", cfg.Database.Host, cfg.Database.Port)
	fmt.Println("Kafka brokers:", cfg.Kafka.Brokers)
}
