CREATE TABLE orders (
   id serial primary key,
   identifier varchar not null unique,
   client integer NOT NULL REFERENCES users(id),
   accrual REAL,
   points_write_off REAL,
   status varchar not null,
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX orders__client__identifier__index
ON orders (client,identifier);