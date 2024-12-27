CREATE TABLE order (
   id serial primary key,
   identifier varchar not null unique,
   client integer NOT NULL REFERENCES user(id),
   accrual integer
   status varchar not null,
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX order__client__identifier__index
ON order (client,identifier);