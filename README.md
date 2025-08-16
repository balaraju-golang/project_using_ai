# project_using_ai
# Banking Microservice

Go microservice for bank accounts and transactions.

## Install

```bash
git clone <repo>
cd banking-service
go mod tidy
```

## Run

```bash
go run cmd/server/main.go
```

## API

POST /accounts
```json
{
  "customer_name": "Ravi Kumar",
  "initial_balance": 10000
}
```

GET /accounts/{id}

POST /transactions/deposit
```json
{
  "account_id": "uuid",
  "amount": 5000
}
```

POST /transactions/withdraw
```json
{
  "account_id": "uuid",
  "amount": 2000
}
```

POST /transactions/transfer
```json
{
  "from_account_id": "uuid",
  "to_account_id": "uuid",
  "amount": 3000
}
```

## Test

```bash
go test ./...
```

## Config

PORT=8080
LOG_LEVEL=info 
