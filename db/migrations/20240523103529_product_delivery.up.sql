CREATE TABLE product_delivery
(
  product_id uuid REFERENCES product (id),
  method_id    uuid REFERENCES delivery_method (id),
  PRIMARY KEY (product_id, method_id)
);