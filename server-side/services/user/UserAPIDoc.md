# User Service API Documentation

## Overview
The User Service API provides endpoints to manage user registration, login, email verification, and user details updates. It interacts with a MySQL database to store user information.

### Base URL
The service is available at `http://localhost:5001`.

## Endpoints

### 1. **GET /api/v1/status**
Check the connection status of the User Service to the database.

#### Response
- **200 OK**: User Service is connected to the database.
- **500 Internal Server Error**: Database connection failed.

#### Example Request:
```
GET /api/v1/status
```

#### Example Response:
```
"User Service connected to the database successfully!"
```

---

### 2. **POST /api/v1/register**
Register a new user.

#### Request Body
- `Name` (string): The name of the user.
- `Email` (string): The email of the user.
- `ContactNo` (string): The contact number of the user.
- `HashedPassword` (string): The password (will be hashed before storing).

#### Response
- **200 OK**: User registered successfully and a verification code sent.
- **400 Bad Request**: Invalid input.
- **409 Conflict**: Email already exists.
- **500 Internal Server Error**: Error during registration.

#### Example Request:
```
POST /api/v1/register
{
  "Name": "John Doe",
  "Email": "john@example.com",
  "ContactNo": "1234567890",
  "HashedPassword": "password123"
}
```

#### Example Response:
```
{
  "message": "User registered successfully. A verification code has been sent to your email."
}
```

---

### 3. **POST /api/v1/verify-email**
Verify the email using the sent verification code.

#### Request Body
- `Email` (string): The email of the user.
- `VerificationCode` (string): The verification code sent to the user's email.

#### Response
- **200 OK**: Email successfully verified.
- **400 Bad Request**: Email already verified.
- **401 Unauthorized**: Invalid verification code.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Error during email verification.

#### Example Request:
```
POST /api/v1/verify-email
{
  "Email": "john@example.com",
  "VerificationCode": "ABCD1234"
}
```

#### Example Response:
```
{
  "message": "Email successfully verified"
}
```

---

### 4. **POST /api/v1/login**
Login a user by email and password.

#### Request Body
- `Email` (string): The email of the user.
- `HashedPassword` (string): The password of the user.

#### Response
- **200 OK**: Login successful.
- **400 Bad Request**: Invalid credentials.
- **401 Unauthorized**: Email not verified.
- **500 Internal Server Error**: Error during login.

#### Example Request:
```
POST /api/v1/login
{
  "Email": "john@example.com",
  "HashedPassword": "password123"
}
```

#### Example Response:
```
{
  "message": "Login successful"
}
```

---

### 5. **GET /api/v1/user/{email}**
Get user details by email.

#### Path Parameters
- `email` (string): The email of the user.

#### Response
- **200 OK**: User details retrieved successfully.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Error querying user details.

#### Example Request:
```
GET /api/v1/user/john@example.com
```

#### Example Response:
```
{
  "userID": 123,
  "Name": "John Doe",
  "Email": "john@example.com",
  "ContactNo": "1234567890",
  "MembershipTier": "Gold",
  "EmailVerified": true
}
```

---

### 6. **PUT /api/v1/user/{email}**
Update user details by email.

#### Path Parameters
- `email` (string): The email of the user.

#### Request Body
- `Name` (string, optional): The updated name of the user.
- `ContactNo` (string, optional): The updated contact number of the user.
- `Password` (string, optional): The updated password (will be hashed before storing).

#### Response
- **200 OK**: User details updated successfully.
- **400 Bad Request**: Invalid input.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Error updating user details.

#### Example Request:
```
PUT /api/v1/user/john@example.com
{
  "Name": "Johnathan Doe",
  "ContactNo": "0987654321",
  "Password": "newpassword123"
}
```

#### Example Response:
```
{
  "message": "User details updated successfully."
}
```

