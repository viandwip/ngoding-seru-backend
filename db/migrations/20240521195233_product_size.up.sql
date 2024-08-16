CREATE TABLE product_size
(
  product_id uuid REFERENCES product (id),
  size_id    uuid REFERENCES size (id),
  price      INTEGER,
  PRIMARY KEY (product_id, size_id)
);
