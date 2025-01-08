CREATE TABLE accounts (
   id serial primary key,
   points REAL,
   withdrawn REAL,
   client integer unique NOT NULL REFERENCES users(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX accounts__client__index
ON accounts (client);