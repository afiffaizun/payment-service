# Payment Service API Documentation

## Base URL
```
http://localhost:8080
```

## Postman Collection

### 1. Health Check
**Method:** GET  
**URL:** `http://localhost:8080/`  
**Description:** Check if the service is running

**Expected Response:**
```
Payment Service is running
```

---

### 2. Transfer Funds
**Method:** POST  
**URL:** `http://localhost:8080/transfer`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
  "sender_id": "user-123",
  "receiver_id": "user-456",
  "amount": 10000,
  "reference": "TRX-20240219-001"
}
```

**Success Response (200 OK):**
```json
{
  "TransactionID": "uuid-generated-id",
  "Reference": "TRX-20240219-001",
  "Amount": 10000,
  "Status": "completed",
  "CreatedAt": "2024-02-19T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body, insufficient balance, same user transfer, or duplicate reference
- `404 Not Found` - Wallet not found

**Postman Tests (Pre-request Script):**
```javascript
// Generate unique reference ID for each request
pm.globals.set("reference_id", "TRX-" + Date.now());
```

**Postman Tests (Tests Tab):**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property("TransactionID");
    pm.expect(jsonData).to.have.property("Reference");
    pm.expect(jsonData).to.have.property("Amount");
    pm.expect(jsonData).to.have.property("Status");
    pm.expect(jsonData).to.have.property("CreatedAt");
});

pm.test("Status is completed", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.Status).to.eql("completed");
});
```

---

### 3. Top Up Wallet
**Method:** POST  
**URL:** `http://localhost:8080/topup`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
  "user_id": "user-123",
  "amount": 50000
}
```

**Success Response (200 OK):**
```json
{
  "UserID": "user-123",
  "Amount": 50000,
  "Balance": 150000
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or invalid amount
- `404 Not Found` - Wallet not found

**Postman Tests (Tests Tab):**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property("UserID");
    pm.expect(jsonData).to.have.property("Amount");
    pm.expect(jsonData).to.have.property("Balance");
});

pm.test("Balance is greater than or equal to amount", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.Balance).to.be.at.least(jsonData.Amount);
});
```

---

### 4. Get Transaction by Reference ID
**Method:** GET  
**URL:** `http://localhost:8080/transaction/{{refId}}`  
**Path Variable:** `refId` - Transaction reference ID

**Example URL:**
```
http://localhost:8080/transaction/TRX-20240219-001
```

**Success Response (200 OK):**
```json
{
  "TransactionID": "uuid-generated-id",
  "Reference": "TRX-20240219-001",
  "Amount": 10000,
  "Status": "completed",
  "CreatedAt": "2024-02-19T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Reference ID is required
- `404 Not Found` - Transaction not found

**Postman Tests (Tests Tab):**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property("TransactionID");
    pm.expect(jsonData).to.have.property("Reference");
    pm.expect(jsonData).to.have.property("Amount");
    pm.expect(jsonData).to.have.property("Status");
    pm.expect(jsonData).to.have.property("CreatedAt");
});

pm.test("Reference matches request", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.Reference).to.eql(pm.variables.get("refId"));
});
```

---

### 5. Get Wallet by User ID
**Method:** GET  
**URL:** `http://localhost:8080/wallet/{{userId}}`  
**Path Variable:** `userId` - User ID

**Example URL:**
```
http://localhost:8080/wallet/user-123
```

**Success Response (200 OK):**
```json
{
  "wallet_id": "wallet-uuid",
  "user_id": "user-123",
  "balance": 150000,
  "version": 5,
  "created_at": "2024-01-15T08:00:00Z",
  "updated_at": "2024-02-19T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - User ID is required
- `404 Not Found` - Wallet not found

**Postman Tests (Tests Tab):**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property("wallet_id");
    pm.expect(jsonData).to.have.property("user_id");
    pm.expect(jsonData).to.have.property("balance");
    pm.expect(jsonData).to.have.property("version");
    pm.expect(jsonData).to.have.property("created_at");
    pm.expect(jsonData).to.have.property("updated_at");
});

pm.test("User ID matches request", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.user_id).to.eql(pm.variables.get("userId"));
});

pm.test("Balance is non-negative", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.balance).to.be.at.least(0);
});
```

---

## Environment Variables

Create a Postman Environment with these variables:

| Variable | Initial Value | Current Value | Description |
|----------|--------------|---------------|-------------|
| base_url | http://localhost:8080 | http://localhost:8080 | Base URL for API |
| sender_id | user-123 | user-123 | Sender user ID for transfers |
| receiver_id | user-456 | user-456 | Receiver user ID for transfers |
| reference_id | TRX-001 | TRX-001 | Transaction reference ID |

---

## Complete Postman Collection JSON

You can import this collection directly into Postman:

```json
{
  "info": {
    "_postman_id": "payment-service-collection",
    "name": "Payment Service API",
    "description": "API collection for Payment Service",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/",
          "host": ["{{base_url}}"],
          "path": [""]
        }
      }
    },
    {
      "name": "Transfer Funds",
      "event": [
        {
          "listen": "prerequest",
          "script": {
            "exec": [
              "pm.globals.set(\"reference_id\", \"TRX-\" + Date.now());"
            ],
            "type": "text/javascript"
          }
        },
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test(\"Status code is 200\", function () {",
              "    pm.response.to.have.status(200);",
              "});",
              "",
              "pm.test(\"Response has required fields\", function () {",
              "    var jsonData = pm.response.json();",
              "    pm.expect(jsonData).to.have.property(\"TransactionID\");",
              "    pm.expect(jsonData).to.have.property(\"Reference\");",
              "    pm.expect(jsonData).to.have.property(\"Amount\");",
              "    pm.expect(jsonData).to.have.property(\"Status\");",
              "    pm.expect(jsonData).to.have.property(\"CreatedAt\");",
              "});"
            ],
            "type": "text/javascript"
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"sender_id\": \"{{sender_id}}\",\n  \"receiver_id\": \"{{receiver_id}}\",\n  \"amount\": 10000,\n  \"reference\": \"{{reference_id}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/transfer",
          "host": ["{{base_url}}"],
          "path": ["transfer"]
        }
      }
    },
    {
      "name": "Top Up Wallet",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test(\"Status code is 200\", function () {",
              "    pm.response.to.have.status(200);",
              "});",
              "",
              "pm.test(\"Response has required fields\", function () {",
              "    var jsonData = pm.response.json();",
              "    pm.expect(jsonData).to.have.property(\"UserID\");",
              "    pm.expect(jsonData).to.have.property(\"Amount\");",
              "    pm.expect(jsonData).to.have.property(\"Balance\");",
              "});"
            ],
            "type": "text/javascript"
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"user_id\": \"{{sender_id}}\",\n  \"amount\": 50000\n}"
        },
        "url": {
          "raw": "{{base_url}}/topup",
          "host": ["{{base_url}}"],
          "path": ["topup"]
        }
      }
    },
    {
      "name": "Get Transaction",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test(\"Status code is 200\", function () {",
              "    pm.response.to.have.status(200);",
              "});",
              "",
              "pm.test(\"Response has required fields\", function () {",
              "    var jsonData = pm.response.json();",
              "    pm.expect(jsonData).to.have.property(\"TransactionID\");",
              "    pm.expect(jsonData).to.have.property(\"Reference\");",
              "    pm.expect(jsonData).to.have.property(\"Amount\");",
              "    pm.expect(jsonData).to.have.property(\"Status\");",
              "    pm.expect(jsonData).to.have.property(\"CreatedAt\");",
              "});"
            ],
            "type": "text/javascript"
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/transaction/{{reference_id}}",
          "host": ["{{base_url}}"],
          "path": ["transaction", "{{reference_id}}"]
        }
      }
    },
    {
      "name": "Get Wallet",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test(\"Status code is 200\", function () {",
              "    pm.response.to.have.status(200);",
              "});",
              "",
              "pm.test(\"Response has required fields\", function () {",
              "    var jsonData = pm.response.json();",
              "    pm.expect(jsonData).to.have.property(\"wallet_id\");",
              "    pm.expect(jsonData).to.have.property(\"user_id\");",
              "    pm.expect(jsonData).to.have.property(\"balance\");",
              "    pm.expect(jsonData).to.have.property(\"version\");",
              "    pm.expect(jsonData).to.have.property(\"created_at\");",
              "    pm.expect(jsonData).to.have.property(\"updated_at\");",
              "});"
            ],
            "type": "text/javascript"
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/wallet/{{sender_id}}",
          "host": ["{{base_url}}"],
          "path": ["wallet", "{{sender_id}}"]
        }
      }
    }
  ]
}
```

---

## Test Scenarios

### Scenario 1: Complete Payment Flow
1. Top up wallet for User A
2. Check User A wallet balance
3. Transfer from User A to User B
4. Get transaction details by reference ID
5. Verify both wallet balances

### Scenario 2: Error Handling
1. Transfer with insufficient balance
2. Transfer to same user
3. Transfer with duplicate reference
4. Get non-existent transaction
5. Get non-existent wallet

### Scenario 3: Concurrent Transfers
1. Create multiple transfers simultaneously
2. Verify atomicity (no partial transfers)
3. Check final balances are consistent

---

## Notes

- All amounts are in the smallest currency unit (e.g., cents)
- Reference IDs must be unique for each transaction
- The service uses optimistic locking for concurrent access
- All endpoints return JSON responses
- Error responses have the format: `{"error": "error message"}`
