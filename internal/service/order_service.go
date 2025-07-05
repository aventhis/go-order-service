package service

import (
    "github.com/aventhis/go-order-service/internal/model"
    "github.com/aventhis/go-order-service/internal/repository"
    "log"
    "sync"
)

// Compile-time interface check
var _ OrderServiceInterface = (*OrderService)(nil)

type OrderService struct {
    repo  repository.OrderRepository
    cache map[string]*model.Order
    mu    sync.RWMutex
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
    return &OrderService{
        repo:  repo,
        cache: make(map[string]*model.Order),
    }
}

// Сервис делает:
// Проверяет кэш: есть ли уже заказ в памяти (map)
// Если нет — идёт в репозиторий (БД)
// Если из БД успешно получили — кладём в кэш
// Возвращаем результат вызывающему коду (например, HTTP-хендлеру)
func (s *OrderService) GetByID(orderUID string) (*model.Order, error) {
    // Сначала проверяем кэш
    s.mu.RLock()
    if order, exists := s.cache[orderUID]; exists {
        s.mu.RUnlock()
        return order, nil
    }
    s.mu.RUnlock()

    // Если нет в кэше, берем из БД
    order, err := s.repo.GetByID(orderUID)
    if err != nil {
        return nil, err
    }

    // Сохраняем в кэш
    s.mu.Lock()
    s.cache[orderUID] = order
    s.mu.Unlock()

    return order, nil
}

// Create сохраняет заказ в БД и кэш
func (s *OrderService) Create(order *model.Order) error {
    if err := s.repo.Create(order); err != nil {
        return err
    }

    s.mu.Lock()
    s.cache[order.OrderUID] = order
    s.mu.Unlock()

    return nil
}

// RestoreCache восстанавливает кэш из БД при старте
func (s *OrderService) RestoreCache() error {
    orders, err := s.repo.GetAll() 
    if err != nil {
        return err
    }

    s.mu.Lock()
    for _, order := range orders {
        s.cache[order.OrderUID] = order
    }
    s.mu.Unlock()

    log.Printf("Restored %d orders to cache", len(orders))
    return nil
}