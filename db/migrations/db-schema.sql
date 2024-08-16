CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users
(
  user_id      uuid                                 DEFAULT gen_random_uuid(),
  email        VARCHAR                     NOT NULL unique,
  phone_number VARCHAR                     NOT NULL unique,
  password     VARCHAR                     NOT NULL,
  role         VARCHAR                     NOT NULL,
  created_at   TIMESTAMP without time zone not null DEFAULT NOW(),
  updated_at   TIMESTAMP without time zone null,
  CONSTRAINT users_pk PRIMARY KEY (user_id)
);


CREATE TABLE profile
(
  profile_id    uuid                                 DEFAULT gen_random_uuid(),
  user_id       uuid                        NOT NULL,
  first_name    VARCHAR                              default '',
  last_name     VARCHAR                              default '',
  display_name  VARCHAR                              default '',
  gender        VARCHAR                              default '',
  address       VARCHAR                              default '',
  birthday      DATE,
  photo_profile TEXT                                 default '',
  created_at    TIMESTAMP without time zone not null DEFAULT NOW(),
  updated_at    TIMESTAMP without time zone null,
  CONSTRAINT profile_pk primary key (profile_id),
  CONSTRAINT fk_profile_users FOREIGN KEY (user_id)
    REFERENCES users (user_id)
    ON DELETE CASCADE
);


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


CREATE TABLE size
(
  id        uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  size_name VARCHAR
);

CREATE TABLE product_size
(
  product_id uuid REFERENCES product (id),
  size_id    uuid REFERENCES size (id),
  price          INTEGER,
  PRIMARY KEY (product_id, size_id)
);

CREATE TABLE favorite
(
  favorite_id uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  product_id  uuid REFERENCES product (id),
  user_id     uuid REFERENCES users (user_id),
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart
(
  id         uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id    uuid REFERENCES users (user_id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_item
(
  id         uuid    NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  cart_id    uuid REFERENCES cart (id),
  product_id uuid REFERENCES product (id),
  size_id    uuid REFERENCES size (id),
  quantity   INTEGER NOT NULL,
  ordered    BOOLEAN,
  created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);
truncate table cart_item;
ALTER TABLE public.cart_item ADD delivery_method_id uuid REFERENCES delivery_method (id);


delete from cart where user_id = 'fed68249-420a-4e4c-b3cc-87ec1738ee5f';


CREATE TABLE cart_order
(
  order_id          uuid           NULL     DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id           uuid REFERENCES users (user_id),
  delivery_method_id uuid REFERENCES delivery_method (id),
  payment_method_id uuid REFERENCES payment_method (id),
  total_price       INTEGER, -- SubTotal
  taxes             DECIMAL(10, 2) NOT NULL DEFAULT 0,
  shipping          INTEGER,
  status            VARCHAR,
  delivery_address  VARCHAR, -- Default ngambil alamat dari profile alamat
  total_amount      INTEGER, -- Total
  created_at        TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
  updated_at        TIMESTAMP               DEFAULT CURRENT_TIMESTAMP
);
DROP TABLE cart_order;
-- SubTotal
CREATE TABLE cart_order_items
(
  order_item_id uuid    NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  order_id      uuid REFERENCES cart_order (order_id),
  user_id       uuid REFERENCES users (user_id),
  product_id    uuid REFERENCES product (id),
  size_id       uuid REFERENCES size (id),
  quantity      INTEGER NOT NULL,
  amount        INTEGER,
  created_at    TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE delivery_method
(
  id          uuid         NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  method_name VARCHAR(100) NOT NULL
);


CREATE TABLE payment_method
(
  id          uuid         NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  method_name VARCHAR(100) NOT NULL
);

CREATE TABLE product_delivery
(
  product_id uuid REFERENCES product (id),
  method_id  uuid REFERENCES delivery_method (id),
  PRIMARY KEY (product_id, method_id)
);

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


DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;



SELECT ci.id                    as cart_item_id,
       ci.product_id,
       ci.size_id,
       ci.quantity,
       ci.created_at,
       ci.updated_at,
       (ci.quantity * ps.price) as amount
FROM cart_item ci
       JOIN product_size ps ON ci.product_id = ps.product_id AND ci.size_id = ps.size_id
WHERE ci.cart_id IN (SELECT id FROM cart WHERE user_id = '8f49f5a3-8a6e-466b-ae51-eef1424e8da1')
  AND ci.ordered = TRUE;
SELECT id         as cart_item_id,
       product_id,
       size_id,
       quantity,
       created_at,
       updated_at,
       (quantity) as amount
FROM cart_item;

SELECT ci.id,
       p.name   as product_name,
       ps.price as product_price,
       ci.quantity,
       ci.created_at,
       ci.updated_at
FROM cart_item ci
       JOIN
     product p ON ci.product_id = p.id
       JOIN
     product_size ps ON ci.product_id = ps.product_id AND ci.size_id = ps.size_id
       JOIN
     cart c ON ci.cart_id = c.id
WHERE c.user_id = 'aabc969d-7e5e-4ea1-9105-35f4da186e50'
  AND ci.ordered = TRUE;

SELECT photo_profile, address, display_name, first_name, last_name, gender, birthday
FROM profile
WHERE user_id = '46350b58-54a3-4a67-a366-a5ee0b7c8266';

select * from cart_order where order_id = 'f7afda74-6bc1-49ed-8f33-84d9b585a8f3';

select total_price, created_at from cart_order where user_id = 'fed68249-420a-4e4c-b3cc-87ec1738ee5f';

SELECT
  SUM(total_price) AS total_price_sum
FROM
  cart_order
WHERE
  created_at >= NOW() - INTERVAL '1 WEEK';

SELECT h.id, h.user_id, h.order_id, h.product_id, p.name as product_name, ps.price as product_price, h.status, h.created_at, h.updated_at
FROM history h
       JOIN product p ON h.product_id = p.id
       JOIN product_size ps ON h.product_id = ps.product_id AND h.size_id = ps.size_id
WHERE h.user_id = '87f938b3-5060-48ca-a4c8-acf67b12f950' AND h.status = 'delivered';