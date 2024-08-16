CREATE TABLE favorite
(
  favorite_id uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  product_id  uuid REFERENCES product (id),
  user_id     uuid REFERENCES users (user_id),
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
