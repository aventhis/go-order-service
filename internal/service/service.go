package service

import "github.com/aventhis/go-order-service/internal/model"

type OrderServiceInterface interface {
    GetByID(orderUID string) (*model.Order, error)
    Create(order *model.Order) error
    RestoreCache() error
}