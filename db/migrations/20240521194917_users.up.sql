CREATE TABLE users
(  
  user_id uuid DEFAULT gen_random_uuid(),
  email VARCHAR NOT NULL unique,
  phone_number VARCHAR NOT NULL unique,
  password VARCHAR NOT NULL,
  role VARCHAR NOT NULL,
	created_at TIMESTAMP without time zone not null DEFAULT NOW(),
	updated_at TIMESTAMP without time zone null,
	CONSTRAINT users_pk PRIMARY KEY (user_id)
);