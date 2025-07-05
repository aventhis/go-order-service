package repository

import (
	"github.com/aventhis/go-order-service/internal/model"
	"github.com/jmoiron/sqlx"
	"log"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(orderUID string) (*model.Order, error)
	GetAll() ([]*model.Order, error)
}

type OrderRepo struct {
    db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepo {
    return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *model.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вставка основной информации о заказе
    _, err = tx.Exec(`
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
        order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
        order.InternalSignature, order.CustomerId, order.DeliveryService,
        order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
    if err != nil {
        return err
    }

	// Вставка информации о доставке
    _, err = tx.Exec(`
        INSERT INTO delivery (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
        order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
        order.Delivery.Region, order.Delivery.Email)
    if err != nil {
        return err
    }

    // Вставка информации об оплате
    _, err = tx.Exec(`
        INSERT INTO payment (
            transaction, order_uid, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
        order.Payment.Transaction, order.OrderUID, order.Payment.RequestID,
        order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
        order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
        order.Payment.GoodsTotal, order.Payment.CustomFee)
    if err != nil {
        return err
    }

	// Вставка информации о товарах
    for _, item := range order.Items {
        _, err = tx.Exec(`
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid,
                name, sale, size, total_price, nm_id, brand, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
            order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
            item.RID, item.Name, item.Sale, item.Size, item.TotalPrice,
            item.NmID, item.Brand, item.Status)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *OrderRepo) GetByID(orderUID string) (*model.Order, error) {
	order := &model.Order{}

	// Получаем основную информацию о заказе
	err := r.db.Get(order, `
		SELECT order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о доставке
	err = r.db.Get(&order.Delivery, `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию об оплате
	err = r.db.Get(&order.Payment, `
		SELECT transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о товарах
	err = r.db.Select(&order.Items, `
		SELECT chrt_id, track_number, price, rid, name, sale,
			size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) GetAll() ([]*model.Order, error) {
    // Получаем все order_uid
    var orders []*model.Order
    err := r.db.Select(&orders, `SELECT * FROM orders`)
    if err != nil {
        return nil, err
    }

    // Для каждого заказа получаем связанные данные
    for _, order := range orders {
        // Получаем информацию о доставке
        err = r.db.Get(&order.Delivery, `
            SELECT name, phone, zip, city, address, region, email
            FROM delivery WHERE order_uid = $1`, order.OrderUID)
        if err != nil {
            return nil, err
        }

        // Получаем информацию об оплате
        err = r.db.Get(&order.Payment, `
            SELECT transaction, request_id, currency, provider, amount,
                payment_dt, bank, delivery_cost, goods_total, custom_fee
            FROM payment WHERE order_uid = $1`, order.OrderUID)
        if err != nil {
            return nil, err
        }

        // Получаем информацию о товарах
        err = r.db.Select(&order.Items, `
            SELECT chrt_id, track_number, price, rid, name, sale,
                size, total_price, nm_id, brand, status
            FROM items WHERE order_uid = $1`, order.OrderUID)
        if err != nil {
            return nil, err
        }
    }

    log.Printf("Found %d orders in database", len(orders))
    return orders, nil
}