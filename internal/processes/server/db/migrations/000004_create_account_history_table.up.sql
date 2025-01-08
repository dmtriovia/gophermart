CREATE TABLE accounts_history (
   id serial primary key,
   points_write_off REAL,
   client_order integer unique NOT NULL REFERENCES orders(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX accounts_history__client_order__index
ON accounts_history (client_order);