package service

import "github.com/aventhis/go-order-service/internal/model"

type OrderService interface {
    GetByID(orderUID string) (*model.Order, error)
}