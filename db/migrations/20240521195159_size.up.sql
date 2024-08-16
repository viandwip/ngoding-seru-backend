CREATE TABLE size
(
  id        uuid NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  size_name VARCHAR
);
INSERT INTO size (size_name)
VALUES ('R'),
       ('L'),
       ('XL');
