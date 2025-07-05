package main

import (
	"fmt"
	"github.com/aventhis/go-order-service/internal/config"
	"github.com/aventhis/go-order-service/internal/db"
	"log"
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

	 // Инициализируем подключение к БД
	database := db.Init(cfg.Database)
	defer database.Close()

	// Применяем миграции
	db.RunMigrations(database, "migrations")

	// Создаем репозиторий
    repo := repository.NewOrderRepository(database)

	// Создаем сервис
	orderService := service.NewOrderService(repo)

	// Создаем HTTP хендлеры
	handler := handlers.NewHandler(orderService)

	// Запускаем сервер
    log.Printf("Starting server on :%d", cfg.Server.Port)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), handler.InitRoutes()); err != nil {
        log.Fatal(err)
    }
}
