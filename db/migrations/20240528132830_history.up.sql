CREATE TABLE history
(
  id             uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id        uuid REFERENCES users (user_id) ON DELETE CASCADE,
  order_id       uuid REFERENCES cart_order (order_id) ON DELETE CASCADE,
  product_id     uuid REFERENCES product (id) ON DELETE CASCADE,
  size_id        uuid REFERENCES size (id) ON DELETE CASCADE,
  status         VARCHAR,
  created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);