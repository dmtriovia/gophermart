CREATE TABLE order (
   id serial primary key,
   identifier varchar not null unique,
   client integer NOT NULL REFERENCES user(id),
   createddate TIMESTAMP default now()
);