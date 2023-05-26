# mini-wallet
Mini wallet service backend application

## Design Document
Please refer to this [notion site](https://boiling-handstand-d11.notion.site/Mini-Wallet-Exercise-11c1786ae8964f3bb2ff0badf109ac2b)

## Setup
1. Install Go
2. Install PostgreSQL
3. Install [Golang Migrate](https://github.com/golang-migrate/migrate)
4. Run command in terminal to migrate database `migrate -database postgres://{username}:{password}@{host}:{port}/mini_wallet?sslmode=disable -path db/migrations up`
5. Start golang application `go run ./cmd/server/...`

## Database Migrations
[Refer to this repository for complete usage](https://github.com/golang-migrate/migrate)

### Create new DB migration
Run command in terminal `migrate create -ext sql -dir db/migrations -seq {migration_name}`

### Migrating Database
Run command in terminal `migrate -database postgres://{username}:{password}@{host}:{port}/mini_wallet?sslmode=disable -path db/migrations up`

### Rollback Database
Run command in terminal `migrate -database postgres://{username}:{password}@{host}:{port}/mini_wallet?sslmode=disable -path db/migrations down`
