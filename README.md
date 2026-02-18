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

This section provides comprehensive testing instructions covering unit tests, integration tests, and end-to-end (E2E) testing.

### Prerequisites

Before running tests, ensure you have:
- Go 1.25 or higher installed
- Docker & Docker Compose running
- PostgreSQL accessible at `localhost:5432`
- Port 8080 available for the application

---

### 1. Unit Testing

Unit tests are located in `internal/usecase/payment_uc_test.go` and use mocked repositories.

**Run all unit tests:**
```bash
go test -v ./internal/usecase/...
```

**Run with coverage report:**
```bash
go test -cover ./internal/usecase/...
```

**Run specific test:**
```bash
go test -v ./internal/usecase/... -run TestTransferFunds
go test -v ./internal/usecase/... -run TestTopUpWallet
```

**Available Unit Tests:**
- ✅ `TestTransferFunds_Success` - Successful money transfer
- ✅ `TestTransferFunds_InsufficientBalance` - Insufficient balance error
- ✅ `TestTransferFunds_InvalidAmount` - Invalid (zero/negative) amount
- ✅ `TestTransferFunds_SameUser` - Transfer to same user error
- ✅ `TestTransferFunds_DuplicateReference` - Duplicate reference ID
- ✅ `TestTopUpWallet_Success` - Successful wallet top-up
- ✅ `TestTopUpWallet_InvalidAmount` - Invalid top-up amount
- ✅ `TestGetWallet_Success` - Get wallet details
- ✅ `TestGetWallet_NotFound` - Wallet not found error

---

### 2. Integration Testing

Integration tests are located in `internal/repository/postgres_repo_test.go` and require a running PostgreSQL database.

**Start test database:**
```bash
docker-compose up -d db
```

**Run integration tests:**
```bash
go test -v ./internal/repository/...
```

**Run specific integration test:**
```bash
go test -v ./internal/repository/... -run TestPostgresRepo_TopUpWallet
```

**Available Integration Tests:**
- ✅ `TestPostgresRepo_TopUpWallet_Success` - Top-up wallet in database
- ✅ `TestPostgresRepo_TopUpWallet_UserNotFound` - Top-up non-existent user
- ✅ `TestPostgresRepo_GetWalletByUserID_Success` - Get wallet from database
- ✅ `TestPostgresRepo_GetWalletByUserID_NotFound` - Get non-existent wallet
- ✅ `TestPostgresRepo_GetTransactionByRef_Success` - Get transaction by reference

---

### 3. End-to-End Testing (Manual)

E2E tests verify the complete flow from HTTP request through all layers to the database.

#### Setup

```bash
# 1. Start database
docker-compose up -d db

# 2. Run application
go run ./cmd/api/main.go

# 3. Verify application is running
curl http://localhost:8080/
# Expected: "Payment Service is running"
```

#### Test Scenarios

##### Test 1: Health Check
```bash
curl http://localhost:8080/
```
**Expected:** `Payment Service is running`

---

##### Test 2: Get Alice's Wallet
```bash
curl http://localhost:8080/wallet/11111111-1111-1111-1111-111111111111
```
**Expected:**
```json
{
  "wallet_id": "...",
  "user_id": "11111111-1111-1111-1111-111111111111",
  "balance": 1000000,
  "version": 0,
  "created_at": "...",
  "updated_at": "..."
}
```

---

##### Test 3: Top-up Alice's Wallet
```bash
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "amount": 50000
  }'
```
**Expected:**
```json
{
  "UserID": "11111111-1111-1111-1111-111111111111",
  "Amount": 50000,
  "Balance": 1050000
}
```

---

##### Test 4: Verify Top-up (Check Balance)
```bash
curl http://localhost:8080/wallet/11111111-1111-1111-1111-111111111111
```
**Expected:** `balance` = 1050000

---

##### Test 5: Transfer Funds (Alice to Bob)
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": 10000,
    "reference": "tx-test-001"
  }'
```
**Expected:**
```json
{
  "TransactionID": "...",
  "Reference": "tx-test-001",
  "Amount": 10000,
  "Status": "completed",
  "CreatedAt": "..."
}
```go test -v ./internal/repository/... -run TestPostgresRepo_TopUpWallet


---

##### Test 6: Verify Transfer (Check Balances)
```bash
# Alice's balance should be 1040000
curl http://localhost:8080/wallet/11111111-1111-1111-1111-111111111111

# Bob's balance should be 60000
curl http://localhost:8080/wallet/22222222-2222-2222-2222-222222222222
```

---

##### Test 7: Get Transaction
```bash
curl http://localhost:8080/transaction/tx-test-001
```
**Expected:** Transaction details with `Reference`: "tx-test-001"

---

##### Test 8: Error - Insufficient Balance
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "22222222-2222-2222-2222-222222222222",
    "receiver_id": "11111111-1111-1111-1111-111111111111",
    "amount": 100000,
    "reference": "tx-fail-001"
  }'
```
**Expected:** `{"error":"insufficient balance"}` (HTTP 400)

---

##### Test 9: Error - Duplicate Reference (Idempotency)
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": 1000,
    "reference": "tx-test-001"
  }'
```
**Expected:** `{"error":"reference ID already exists"}` (HTTP 400)

---

##### Test 10: Error - Invalid Amount (Negative)
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "11111111-1111-1111-1111-111111111111",
    "receiver_id": "22222222-2222-2222-2222-222222222222",
    "amount": -100,
    "reference": "tx-fail-002"
  }'
```
**Expected:** `{"error":"amount must be greater than zero"}` (HTTP 400)

---

##### Test 11: Error - Wallet Not Found
```bash
curl http://localhost:8080/wallet/99999999-9999-9999-9999-999999999999
```
**Expected:** `{"error":"wallet not found"}` (HTTP 404)

---

##### Test 12: Error - Invalid Top-up Amount
```bash
curl -X POST http://localhost:8080/topup \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "amount": 0
  }'
```
**Expected:** `{"error":"amount must be greater than zero"}` (HTTP 400)

---

### 4. Test Scenarios Summary

| # | Test Case | Endpoint | Method | Expected Result |
|---|-----------|----------|--------|-----------------|
| 1 | Health Check | `/` | GET | `Payment Service is running` |
| 2 | Get Wallet | `/wallet/{userId}` | GET | Wallet details with timestamps |
| 3 | Top-up Wallet | `/topup` | POST | Balance increased |
| 4 | Verify Top-up | `/wallet/{userId}` | GET | Updated balance |
| 5 | Transfer Funds | `/transfer` | POST | Transaction completed |
| 6 | Verify Transfer | `/wallet/{userId}` | GET | Both balances updated |
| 7 | Get Transaction | `/transaction/{refId}` | GET | Transaction details |
| 8 | Insufficient Balance | `/transfer` | POST | `error: insufficient balance` (400) |
| 9 | Duplicate Reference | `/transfer` | POST | `error: reference ID already exists` (400) |
| 10 | Invalid Amount | `/transfer` | POST | `error: amount must be greater than zero` (400) |
| 11 | Wallet Not Found | `/wallet/{userId}` | GET | `error: wallet not found` (404) |
| 12 | Invalid Top-up | `/topup` | POST | `error: amount must be greater than zero` (400) |

---

### 5. Postman Collection Setup

#### Collection Structure

Create a Collection called **"Payment Service"** with folders:

**Folder: Health & Info**
- Health Check (GET /)

**Folder: Wallet Operations**
- Get Wallet (GET /wallet/{userId})
- Top-up Wallet (POST /topup)

**Folder: Transfer Operations**
- Transfer Funds (POST /transfer)
- Get Transaction (GET /transaction/{refId})

**Folder: Error Cases**
- Insufficient Balance (POST /transfer)
- Duplicate Reference (POST /transfer)
- Invalid Amount (POST /transfer)
- Wallet Not Found (GET /wallet/{userId})

#### Environment Variables

Create an Environment with these variables:
```
BASE_URL: http://localhost:8080
ALICE_USER_ID: 11111111-1111-1111-1111-111111111111
BOB_USER_ID: 22222222-2222-2222-2222-222222222222
```

#### Example Request Configuration

**Get Wallet:**
- **Method:** GET
- **URL:** `{{BASE_URL}}/wallet/{{ALICE_USER_ID}}`
- **Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response has wallet_id", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('wallet_id');
    pm.expect(jsonData).to.have.property('balance');
});
```

**Transfer Funds:**
- **Method:** POST
- **URL:** `{{BASE_URL}}/transfer`
- **Headers:** `Content-Type: application/json`
- **Body:**
```json
{
  "sender_id": "{{ALICE_USER_ID}}",
  "receiver_id": "{{BOB_USER_ID}}",
  "amount": 5000,
  "reference": "tx-{{$timestamp}}"
}
```

---

### 6. Troubleshooting Tests

#### Problem: Port 8080 already in use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Problem: Database connection refused
```bash
# Check if database container is running
docker-compose ps

# Restart database
docker-compose restart db

# Check logs
docker-compose logs db
```

#### Problem: Tests failing with "no rows in result set"
Database seed data might be missing. Reset the database:
```bash
docker-compose down -v
docker-compose up -d db
```

#### Problem: Unit tests fail with mock errors
Ensure all dependencies are installed:
```bash
go mod tidy
go mod download
```

#### Problem: Permission denied when running test script
```bash
chmod +x run-tests.sh
```

#### Problem: PostgreSQL authentication failed
Check environment variables match docker-compose.yml:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=user_payment
export DB_PASSWORD=pass_payment
export DB_NAME=db_payment
```

---

### 7. Quick Test Checklist

Before committing changes, run through this checklist:

- [ ] **Unit Tests:** `go test ./internal/usecase/...` - All pass
- [ ] **Integration Tests:** `go test ./internal/repository/...` - All pass
- [ ] **Build:** `go build ./cmd/api/main.go` - No errors
- [ ] **Health Check:** `curl http://localhost:8080/` - Returns "Payment Service is running"
- [ ] **Get Wallet:** Returns wallet with correct balance and timestamps
- [ ] **Top-up:** Increases balance correctly
- [ ] **Transfer:** Completes successfully and updates both balances
- [ ] **Error Handling:** All error cases return appropriate status codes and messages
- [ ] **Idempotency:** Duplicate reference returns error

---

### 8. Reset Test Data

To reset database to initial state with dummy data:

```bash
# Stop and remove containers
docker-compose down -v

# Start fresh
docker-compose up -d db

# Wait for database to initialize
sleep 5

# Verify tables and data
docker exec payment-service-db-1 psql -U user_payment -d db_payment -c "SELECT * FROM wallets;"
```

**Initial Dummy Data:**
- **Alice:** `11111111-1111-1111-1111-111111111111` - Balance: 1,000,000
- **Bob:** `22222222-2222-2222-2222-222222222222` - Balance: 50,000

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