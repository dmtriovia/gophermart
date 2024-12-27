CREATE TABLE order (
   id serial primary key,
   identifier varchar not null unique,
   client integer NOT NULL REFERENCES user(id),
   status varchar not null,
   createddate TIMESTAMP default now()
);