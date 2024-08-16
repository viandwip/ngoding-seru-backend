CREATE TABLE cart_order
(
  order_id           uuid NULL     DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id            uuid REFERENCES users (user_id),
  delivery_method_id uuid REFERENCES delivery_method (id),
  payment_method_id  uuid REFERENCES payment_method (id),
  total_price        INTEGER, -- SubTotal
  taxes              DECIMAL(10, 2) NOT NULL DEFAULT 0,
  shipping           INTEGER,
  status             VARCHAR,
  delivery_address   VARCHAR, -- Default ngambil alamat dari profile alamat
  total_amount       INTEGER, -- Total
  created_at         TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
  updated_at         TIMESTAMP               DEFAULT CURRENT_TIMESTAMP
);
