CREATE TABLE users (
   id serial primary key,
   login varchar not null unique,
   password varchar not null,
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX users__login__index
ON users (login);