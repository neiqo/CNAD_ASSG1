# Billing Service API Documentation

## Overview
The Billing Service API provides endpoints to manage payment processing, update payment statuses, and send receipts. It interacts with a database for storing payment records and integrates with other services for promotions and booking updates.

### Base URL
The service is available at `http://localhost:5004`.

## Endpoints

### 1. **GET /api/v1/status**
Check the connection status of the Billing Service to the database.

#### Response
- **200 OK**: Billing Service is connected to the database.
- **500 Internal Server Error**: Database connection failed.

#### Example Request:
```
GET /api/v1/status
```

#### Example Response:
```
"Billing Service connected to the database successfully!"
```

---

### 2. **POST /api/v1/payments**
Create a new payment record.

#### Request Body
- `userID` (integer): The ID of the user making the payment.
- `bookingID` (integer): The booking ID associated with the payment.
- `amount` (float): The total payment amount.
- `promotionID` (integer): The ID of the promotion applied to the payment (optional).

#### Response
- **201 Created**: Payment created successfully.
- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Failed to create payment.

#### Example Request:
```
POST /api/v1/payments
{
  "userID": 123,
  "bookingID": 456,
  "amount": 100.0,
  "promotionID": 789
}
```

#### Example Response:
```
{
  "message": "Payment created and is pending!",
  "paymentID": 101,
  "finalAmount": 80.0
}
```

---

### 3. **PUT /api/v1/payments/{paymentID}**
Update the status of a payment.

#### Path Parameters
- `paymentID` (integer): The ID of the payment to update.

#### Request Body
- `status` (string): The new status of the payment. Example values: `"Pending"`, `"Successful"`, `"Failed"`.

#### Response
- **200 OK**: Payment status updated successfully.
- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Failed to update payment status.

#### Example Request:
```
PUT /api/v1/payments/101
{
  "status": "Successful"
}
```

#### Example Response:
```
{
  "message": "Payment status updated successfully"
}
```

---

### 4. **GET /api/v1/payments/{paymentID}**
Retrieve payment details by payment ID.

#### Path Parameters
- `paymentID` (integer): The ID of the payment to retrieve.

#### Response
- **200 OK**: Payment details retrieved successfully.
- **404 Not Found**: Payment not found.
- **500 Internal Server Error**: Failed to retrieve payment details.

#### Example Request:
```
GET /api/v1/payments/101
```

#### Example Response:
```
{
  "paymentID": 101,
  "userID": 123,
  "bookingID": 456,
  "status": "Pending",
  "promotionID": 789,
  "amount": 80.0,
  "finalAmount": 80.0
}
```
