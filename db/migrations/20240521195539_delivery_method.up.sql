CREATE TABLE delivery_method
(
  id     uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  method_name   VARCHAR(100) NOT NULL
);
INSERT INTO delivery_method (method_name)
VALUES ('Dine In'),
       ('Pick Up'),
       ('Door Delivery');