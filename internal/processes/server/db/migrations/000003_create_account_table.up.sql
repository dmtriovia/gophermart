CREATE TABLE account (
   id serial primary key,
   points REAL,
   withdrawn REAL,
   client integer unique NOT NULL REFERENCES user(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX account__client__index
ON account (client);