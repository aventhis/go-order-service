# Go Order Service

Демонстрационный микросервис на Go для обработки заказов с использованием Kafka, PostgreSQL и in-memory кеша.

## Описание

Сервис получает данные о заказах из Kafka, сохраняет их в PostgreSQL и кэширует в памяти для быстрого доступа. Есть HTTP API и простой веб-интерфейс для просмотра заказа по ID.

---

## Быстрый старт

## ⚙️ Настройка переменных окружения

1. Скопируйте `.env.example` в `.env`:
```bash
cp .env.example .env

### 1. Клонируйте репозиторий

```bash
git clone git@github.com:aventhis/go-order-service.git
cd go-order-service
```

### 2. Запустите все сервисы через Docker Compose

```bash
docker compose up -d
```

> **Важно:**  
> Убедитесь, что Docker Engine запущен и порты 5432/5433, 9092, 8081 свободны.

### 3. Проверьте статус сервисов

```bash
docker compose ps
```

---

## Проверка работы

### 1. Проверить, что заказы поступают в базу

```bash
docker exec -it go-order-service-db-1 psql -U postgres -d orders -c "SELECT COUNT(*) FROM orders;"
```

### 2. Посмотреть последние заказы

```bash
docker exec -it go-order-service-db-1 psql -U postgres -d orders -c "SELECT order_uid, track_number, date_created FROM orders ORDER BY date_created DESC LIMIT 5;"
```

### 3. Получить заказ через HTTP API

```bash
curl -X GET http://localhost:8081/order/<order_uid>
```
_Подставьте реальный order_uid из предыдущего шага._

### 4. Проверить работу кеша

```bash
time curl -X GET http://localhost:8081/order/<order_uid>
time curl -X GET http://localhost:8081/order/<order_uid>
```

### 5. Проверить обработку несуществующего заказа

```bash
curl -X GET http://localhost:8081/order/does-not-exist
```

### 6. Перезапустить сервис и убедиться, что кеш восстанавливается

```bash
docker compose restart order-service
curl -X GET http://localhost:8081/order/<order_uid>
```

### 7. Открыть веб-интерфейс

Перейдите в браузере по адресу:  
[http://localhost:8081](http://localhost:8081)

---

## Структура проекта

- `cmd/` — точка входа приложения и продюсера
- `internal/model/` — модели данных
- `internal/repository/` — работа с БД
- `internal/delivery/` — HTTP и Kafka обработчики
- `migrations/` — SQL-миграции для создания таблиц
- `Dockerfile`, `docker-compose.yml` — инфраструктура

---

## Остановка и очистка

```bash
docker compose down -v
```

---

## Примечания

- Для работы требуется установленный Docker Desktop.
- Все переменные окружения уже прописаны в docker-compose.yml.
- Для тестирования Kafka используется встроенный продюсер (`order-producer`).

---


