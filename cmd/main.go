package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/aventhis/go-order-service/internal/config"
    "github.com/aventhis/go-order-service/internal/db"
    "github.com/aventhis/go-order-service/internal/repository"
    "github.com/aventhis/go-order-service/internal/service"
    httphandler "github.com/aventhis/go-order-service/internal/delivery/http"
    "github.com/aventhis/go-order-service/internal/delivery/kafka"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("App running on port:", cfg.Server.Port)
	fmt.Println("Connected to DB:", cfg.Database.Host, cfg.Database.Port)
	fmt.Println("Kafka brokers:", cfg.Kafka.Brokers)

	 // Иницsиализируем подключение к БД
	database := db.Init(cfg.Database)
	defer database.Close()

	// Применяем миграции
	db.RunMigrations(database, "migrations")

	// Создаем репозиторий
    repo := repository.NewOrderRepository(database)

	// Создаем сервис
	orderService := service.NewOrderService(repo)

	// Восстанавливаем кэш
    if err := orderService.RestoreCache(); err != nil {
        log.Printf("Failed to restore cache: %v", err)
    }

	// Создаем и запускаем Kafka consumer
    consumer := kafka.NewConsumer(cfg.Kafka, orderService)
    if err := consumer.Start(); err != nil {
        log.Printf("Failed to start Kafka consumer: %v", err)
    }
	defer consumer.Close()

	// Создаем HTTP хендлеры
	handler := httphandler.NewHandler(orderService)

	// Запускаем сервер
    log.Printf("Starting server on :%d", cfg.Server.Port)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), handler.InitRoutes()); err != nil {
        log.Fatal(err)
    }
}
