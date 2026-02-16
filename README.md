# Payment Service

A Go backend service for money transfers using Clean Architecture, PostgreSQL, and Docker.

## Prerequisites

- Go 1.25 or higher
- Docker & Docker Compose
- PostgreSQL (optional, if not using Docker)

## Clone & Setup

```bash
git clone <repository-url>
cd payment-service
```

## Start Database

```bash
docker compose up -d
```

The database will be available at `localhost:5432` with:
- Database: `db_payment`
- User: `user_payment`
- Password: `pass_payment`

## Run Application

```bash
go mod tidy
go run ./cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### POST /transfer
Transfer funds between users.

**Request:**
```json
{
  "sender_id": "11111111-1111-1111-1111-111111111111",
  "receiver_id": "22222222-2222-2222-2222-222222222222",
  "amount": 5000,
  "reference": "tx-001"
}
```

**Response:**
```json
{
  "TransactionID": "539f9494-dda7-49f9-b3d9-7947075d513d",
  "Reference": "tx-001",
  "Amount": 5000,
  "Status": "completed",
  "CreatedAt": "2026-02-16T13:36:02.772427148+07:00"
}
```

### GET /transaction/{refId}
Get transaction details by reference ID.

**Response:**
```json
{
  "TransactionID": "539f9494-dda7-49f9-b3d9-7947075d513d",
  "Reference": "tx-001",
  "Amount": 5000,
  "Status": "completed",
  "CreatedAt": "2026-02-16T13:36:02.772427148+07:00"
}
```

## Testing

### Test 1: Successful Transfer
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": 5000,
    "reference": "tx-001"
  }'
```

### Test 2: Get Transaction
```bash
curl http://localhost:8080/transaction/tx-001
```

### Test 3: Insufficient Balance
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "22222222-2222-2222-2222-222222222222",
    "receiver_id": "11111111-1111-1111-1111-111111111111",
    "amount": 1000000,
    "reference": "tx-002"
  }'
```
Expected: `"error":"insufficient balance"`

### Test 4: Duplicate Reference (Idempotency)
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": 1000,
    "reference": "tx-001"
  }'
```
Expected: `"error":"reference ID already exists"`

### Test 5: Invalid Amount
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": -100,
    "reference": "tx-003"
  }'
```
Expected: `"error":"amount must be greater than zero"`

## Dummy Data

The database comes with pre-seeded data for testing:

| User | User ID | Wallet Balance |
|------|---------|----------------|
| Alice | `11111111-1111-1111-1111-111111111111` | 1,000,000 |
| Bob | `22222222-2222-2222-2222-222222222222` | 50,000 |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | Database host |
| DB_PORT | 5432 | Database port |
| DB_USER | user_payment | Database user |
| DB_PASSWORD | pass_payment | Database password |
| DB_NAME | db_payment | Database name |

## Architecture

```
payment-service/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── domain/               # Entities & interfaces
│   │   └── wallet.go
│   ├── usecase/              # Business logic
│   │   └── payment_uc.go
│   ├── repository/           # Database implementation
│   │   └── postgres_repo.go
│   └── delivery/             # HTTP handlers
│       └── http_handler.go
├── migrations/               # SQL schemas
│   └── 0000001_init_schema.up.sql
├── docker-compose.yml
├── go.mod
└── README.md
```

## Features

- ✅ Money transfer between users
- ✅ Database transactions with row locking (SELECT FOR UPDATE)
- ✅ Idempotency using reference IDs
- ✅ Balance validation (no negative balance)
- ✅ Clean Architecture (Domain → Usecase → Repository → Delivery)
