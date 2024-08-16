CREATE TABLE product
(
  id             uuid          NULL     DEFAULT gen_random_uuid() PRIMARY KEY,
  name           VARCHAR,
  description    VARCHAR,
  image_url      VARCHAR,
  category       VARCHAR,
  is_available   BOOLEAN,
  delivery_start VARCHAR,
  delivery_end   VARCHAR,
  created_at     TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);
