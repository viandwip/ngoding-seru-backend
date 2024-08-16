CREATE TABLE payment_method
(
  id   uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  method_name VARCHAR(100) NOT NULL
);

INSERT INTO payment_method (method_name)
VALUES ('Card'),
       ('Bank account'),
       ('Cash on delivery');
