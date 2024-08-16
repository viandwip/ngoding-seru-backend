CREATE TABLE cart_item
(
  id         uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  cart_id    uuid REFERENCES cart (id),
  product_id uuid REFERENCES product (id),
  size_id    uuid REFERENCES size (id),
  quantity   INTEGER NOT NULL,
  ordered    BOOLEAN,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
