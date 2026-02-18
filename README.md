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
# Troubleshooting
```
lsof -i :8080
kill -9 <PID>
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

### POST /topup
Top-up a user's wallet.

**Request:**
```json
{
  "user_id": "11111111-1111-1111-1111-111111111111",
  "amount": 10000
}
```

**Response:**
```json
{
  "UserID": "11111111-1111-1111-1111-111111111111",
  "Amount": 10000,
  "Balance": 1010000
}
```

### GET /wallet/{userId}
Get wallet details by user ID.

**Response:**
```json
{
  "wallet_id": "aaaaaaaa-1111-1111-1111-111111111111",
  "user_id": "11111111-1111-1111-1111-111111111111",
  "balance": 1010000,
  "version": 1,
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-15T14:30:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "wallet not found"
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

### Test 6: Top-up Wallet
```bash
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "amount": 10000
  }'
```
Expected: `{"UserID":"11111111-1111-1111-1111-111111111111","Amount":10000,"Balance":...}`

### Test 7: Get Wallet
```bash
curl http://localhost:8080/wallet/11111111-1111-1111-1111-111111111111
```
Expected: `{"wallet_id":"...","user_id":"11111111-1111-1111-1111-111111111111","balance":...,"created_at":"...","updated_at":"..."}`

### Test 8: Get Wallet - Not Found
```bash
curl http://localhost:8080/wallet/99999999-9999-9999-9999-999999999999
```
Expected: `{"error":"wallet not found"}`

## Testing with Postman

### 1. Start the Services

**Option A: Using Docker (Recommended)**
```bash
docker compose up -d
```

**Option B: Manual**
```bash
docker compose up -d db      # Start database only
go run ./cmd/api/main.go     # Run application
```

### 2. Import to Postman

Create a new Collection called "Payment Service" with these requests:

---

#### Request 1: Health Check
- **Method:** GET
- **URL:** `http://localhost:8080/`
- **Expected:** `Payment Service is running`

---

#### Request 2: Transfer Funds
- **Method:** POST
- **URL:** `http://localhost:8080/transfer`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
```json
{
  "sender_id": "11111111-1111-1111-1111-111111111111",
  "receiver_id": "22222222-2222-2222-2222-222222222222",
  "amount": 5000,
  "reference": "tx-001"
}
```
- **Expected Response:**
```json
{
  "TransactionID": "...",
  "Reference": "tx-001",
  "Amount": 5000,
  "Status": "completed",
  "CreatedAt": "..."
}
```

---

#### Request 3: Get Transaction
- **Method:** GET
- **URL:** `http://localhost:8080/transaction/tx-001`

---

#### Request 4: Top-up Wallet
- **Method:** POST
- **URL:** `http://localhost:8080/topup`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
```json
{
  "user_id": "11111111-1111-1111-1111-111111111111",
  "amount": 10000
}
```
- **Expected Response:**
```json
{
  "UserID": "11111111-1111-1111-1111-111111111111",
  "Amount": 10000,
  "Balance": "..."
}
```

---

#### Request 5: Get Wallet
- **Method:** GET
- **URL:** `http://localhost:8080/wallet/11111111-1111-1111-1111-111111111111`
- **Expected Response:**
```json
{
  "wallet_id": "aaaaaaaa-1111-1111-1111-111111111111",
  "user_id": "11111111-1111-1111-1111-111111111111",
  "balance": 1010000,
  "version": 1,
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-15T14:30:00Z"
}
```

---

### 3. Test Scenarios

| Scenario | Expected Result |
|----------|-----------------|
| **Test 1:** Transfer with insufficient balance (Bob has 50k, try to send 100k) | `"error":"insufficient balance"` |
| **Test 2:** Duplicate reference ID (use tx-001 again) | `"error":"reference ID already exists"` |
| **Test 3:** Negative amount (for transfer) | `"error":"amount must be greater than zero"` |
| **Test 4:** Negative amount (for top-up) | `"error":"amount must be greater than zero"` |
| **Test 5:** Get wallet for non-existent user | `"error":"wallet not found"` |

---

### 4. Dummy Data Available

- **Alice:** `11111111-1111-1111-1111-111111111111` (Balance: 1,000,000)
- **Bob:** `22222222-2222-2222-2222-222222222222` (Balance: 50,000)

---

### 5. Postman Collection Export (Optional)

You can create a Postman collection and export it for sharing with your team.

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
- ✅ Wallet top-up functionality
- ✅ Get wallet details with timestamps
- ✅ Clean Architecture (Domain → Usecase → Repository → Delivery)

## Troubleshooting

```
lsof -i :8080
kill -9 <PID>
```

## GUI Postgres
Connect with:
- Host: localhost
- Port: 5432
- Database: db_payment
- User: user_payment
- Password: pass_payment

### Cek isi di docker
- docker exec payment-service-db-1 psql -U user_payment -d db_payment -c "\dt"

---
D. End-to-End Testing
1. Start Services:
docker-compose up -d db
go run cmd/api/main.go
2. Test Scenarios:
# Test 1: TopUp berhasil
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "amount": 50000
  }'
# Expected: {"user_id":"...","new_balance":1050000}
# Test 2: TopUp dengan amount 0
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{"user_id":"11111111-1111-1111-1111-111111111111","amount":0}'
# Expected: 400 Bad Request, "amount must be greater than zero"
# Test 3: TopUp user tidak ada
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{"user_id":"99999999-9999-9999-9999-999999999999","amount":1000}'
# Expected: 404 Not Found