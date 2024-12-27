CREATE TABLE user (
   id serial primary key,
   login varchar not null unique,
   password varchar not null,
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX order__login__index
ON user (login);