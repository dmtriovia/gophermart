module github.com/dmitrovia/gophermart

go 1.22.0

require (
	github.com/golang-migrate/migrate/v4 v4.18.1
	go.uber.org/zap v1.27.0
)

require github.com/golang-jwt/jwt/v4 v4.5.1 // direct

require github.com/gorilla/mux v1.8.1 // direct

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // direct
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.27.0 // direct
	golang.org/x/text v0.18.0 // indirect
)
