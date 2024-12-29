CREATE TABLE account_history (
   id serial primary key,
   points_write_off REAL,
   order integer unique NOT NULL REFERENCES order(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX account_history__order__index
ON account_history (order);