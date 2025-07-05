DOCKER_COMPOSE = docker-compose
DB_URL = postgres://order_app:secret@localhost:5432/order_service?sslmode=disable

# Запуск приложения локально
run:
	go run ./cmd/main.go

# Поднять docker-compose
docker-up:
	$(DOCKER_COMPOSE) up -d

# Применить миграции
migrate-up:
	migrate -database "$(DB_URL)" -path ./migrations up

# Откатить миграции
migrate-down:
	migrate -database "$(DB_URL)" -path ./migrations down