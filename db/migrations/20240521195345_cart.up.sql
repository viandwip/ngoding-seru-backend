CREATE TABLE cart
(
  id           uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id      uuid REFERENCES users (user_id),
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);