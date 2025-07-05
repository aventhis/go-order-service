package main

import (
    "encoding/json"
    "fmt"
    "github.com/Shopify/sarama"
    "github.com/joho/godotenv"
    "log"
    "os"
    "time"
    "github.com/aventhis/go-order-service/internal/model"
)

// Генерация заказа с одинаковыми данными, но уникальным ID
func generateOrder() *model.Order {
    return &model.Order{
        OrderUID:    fmt.Sprintf("test-order-%d", time.Now().UnixNano()), // Уникальный ID на основе времени
        TrackNumber: "WBILMTESTTRACK",
        Entry:      "WBIL",
        Delivery: model.Delivery{
            Name:    "Test Testov",
            Phone:   "+9720000000",
            Zip:     "2639809",
            City:    "Kiryat Mozkin",
            Address: "Ploshad Mira 15",
            Region:  "Kraiot",
            Email:   "test@gmail.com",
        },
        Payment: model.Payment{
            Transaction:  "b563feb7b2b84b6test",
            RequestID:    "",
            Currency:    "USD",
            Provider:    "wbpay",
            Amount:      1817,
            PaymentDT:   1637907727,
            Bank:        "alpha",
            DeliveryCost: 1500,
            GoodsTotal:   317,
            CustomFee:    0,
        },
        Items: []model.Item{
            {
                ChrtID:      9934930,
                TrackNumber: "WBILMTESTTRACK",
                Price:      453,
                RID:        "ab4219087a764ae0btest",
                Name:       "Mascaras",
                Sale:      30,
                Size:      "0",
                TotalPrice: 317,
                NmID:      2389212,
                Brand:     "Vivienne Sabo",
                Status:    202,
            },
        },
        Locale:            "en",
        InternalSignature: "",
        CustomerId:        "test",
        DeliveryService:   "meest",
        Shardkey:         "9",
        SmId:             99,
        DateCreated:      time.Now(),
        OofShard:         "1",
    }
}

func main() {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }

    // Получаем настройки из env
    kafkaBroker := os.Getenv("KAFKA_BROKERS")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }
    
    kafkaTopic := os.Getenv("KAFKA_TOPIC")
    if kafkaTopic == "" {
        kafkaTopic = "orders"
    }

    // Интервал между отправкой заказов (в секундах)
    interval := os.Getenv("PRODUCER_INTERVAL")
    if interval == "" {
        interval = "5" // по умолчанию 5 секунд
    }
    intervalInt := 5
    fmt.Sscanf(interval, "%d", &intervalInt)

    config := sarama.NewConfig()
    config.Producer.Return.Successes = true

    producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    log.Printf("Starting producer with %d second interval\n", intervalInt)
    
    // Бесконечный цикл генерации и отправки заказов
    for {
        order := generateOrder()
        orderBytes, err := json.Marshal(order)
        if err != nil {
            log.Printf("Error marshaling order: %v", err)
            continue
        }

        msg := &sarama.ProducerMessage{
            Topic: kafkaTopic,
            Value: sarama.StringEncoder(orderBytes),
        }

        partition, offset, err := producer.SendMessage(msg)
        if err != nil {
            log.Printf("Error sending message: %v", err)
            continue
        }

        log.Printf("Sent order %s to partition %d at offset %d\n", order.OrderUID, partition, offset)

        // Ждем перед отправкой следующего заказа
        time.Sleep(time.Duration(intervalInt) * time.Second)
    }
}