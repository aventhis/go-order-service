FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum, устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd

# Этап 2: минимальный образ
FROM alpine:latest

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/server .

EXPOSE 8081

# Запуск
CMD ["./server"]