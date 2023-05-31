# Graduate work

# To start and enter db:
docker-compose -f db.yaml up -d

psql postgresql://debts:debts@localhost:5432

# To make and run migrations:
https://github.com/golang-migrate/migrate

migrate create -ext sql -dir ./schema -seq init

migrate -path ./schema -database 'postgresql://debts:debts@localhost:5432/debts?sslmode=disable' up
