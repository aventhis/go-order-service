package main

import (
	"fmt"
	"github.com/aventhis/go-order-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка чтения конфиг файла: %v", err)
	}

	fmt.Println("App running on port:", cfg.Server.Port)
	fmt.Println("Connected to DB:", cfg.Database.Host, cfg.Database.Port)
	fmt.Println("Kafka brokers:", cfg.Kafka.Brokers)
}
