CREATE TABLE user (
   id serial primary key,
   login varchar not null,
   password varchar not null,
   createddate TIMESTAMP default now()
);