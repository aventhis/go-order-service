-- СОЗДАНИЕ ТАБЛИЦЫ ЗАКАЗОВ
CREATE TABLE orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard TEXT
);

-- Таблица доставки (1:1 с orders)
CREATE TABLE delivery (
    order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid),
    name      TEXT,
    phone     TEXT,
    zip       TEXT,
    city      TEXT,
    address   TEXT,
    region    TEXT,
    email     TEXT
);

-- Таблица платежей (1:1 с orders)
CREATE TABLE payment (
    transaction     TEXT PRIMARY KEY,
    order_uid       TEXT UNIQUE REFERENCES orders(order_uid),
    request_id      TEXT,
    currency        TEXT,
    provider        TEXT,
    amount          INT,
    payment_dt      BIGINT,
    bank            TEXT,
    delivery_cost   INT,
    goods_total     INT,
    custom_fee      INT
);

-- Таблица товаров (многие к одному — order_uid)
CREATE TABLE items (
    id            SERIAL PRIMARY KEY,
    order_uid     TEXT REFERENCES orders(order_uid),
    chrt_id       INT,
    track_number  TEXT,
    price         INT,
    rid           TEXT,
    name          TEXT,
    sale          INT,
    size          TEXT,
    total_price   INT,
    nm_id         INT,
    brand         TEXT,
    status        INT
);