CREATE TABLE profile
(
  profile_id    uuid    DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL,
  first_name    VARCHAR default '',
  last_name     VARCHAR default '',
  display_name  VARCHAR default '',
  gender        VARCHAR default '',
  address       VARCHAR default '',
  birthday      DATE,
  photo_profile TEXT    default '',
  created_at    TIMESTAMP without time zone not null DEFAULT NOW(),
  updated_at    TIMESTAMP without time zone null,
  CONSTRAINT profile_pk primary key (profile_id),
  CONSTRAINT fk_profile_users FOREIGN KEY (user_id)
    REFERENCES users (user_id)
    ON DELETE CASCADE
);
