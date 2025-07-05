package kafka

import (
    "encoding/json"
    "log"
    "github.com/Shopify/sarama"
    "github.com/aventhis/go-order-service/internal/config"
    "github.com/aventhis/go-order-service/internal/model"
)

type OrderHandler interface {
    Create(order *model.Order) error
}

type Consumer struct {
    handler OrderHandler
    cfg     config.KafkaConfig
	consumer  sarama.Consumer
    partition sarama.PartitionConsumer
}

func NewConsumer(cfg config.KafkaConfig, handler OrderHandler) *Consumer {
    return &Consumer{
        handler: handler,
        cfg:     cfg,
    }
}

func (c *Consumer) Start() error {
    consumer, err := sarama.NewConsumer([]string{c.cfg.Brokers}, nil)
    if err != nil {
        return err
    }
	c.consumer = consumer

    partition, err := consumer.ConsumePartition(c.cfg.Topic, 0, sarama.OffsetNewest)
    if err != nil {
        return err
    }
	c.partition = partition

    go func() {
        for msg := range partition.Messages() {
            var order model.Order
            if err := json.Unmarshal(msg.Value, &order); err != nil {
                log.Printf("Error parsing order: %v", err)
                continue
            }

            if err := c.handler.Create(&order); err != nil {
                log.Printf("Error saving order: %v", err)
            }
        }
    }()

    return nil
}


func (c *Consumer) Close() error {
    if c.partition != nil {
        if err := c.partition.Close(); err != nil {
            return err
        }
    }
    if c.consumer != nil {
        if err := c.consumer.Close(); err != nil {
            return err
        }
    }
    return nil
}