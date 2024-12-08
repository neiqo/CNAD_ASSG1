

## Vehicle Service API Documentation

### Base URL
```
http://localhost:5000/api/v1
```

---

### **Endpoints**

#### **1. GET /status**
- **Description**: Checks if the Vehicle Service is connected to the database.
- **Method**: `GET`
- **URL**: `/status`
- **Response**:
  - **200 OK**: `"Vehicle Service connected to the database successfully!"`
  - **500 Internal Server Error**: `"Error: Vehicle Service failed to connect to the database"`

#### **2. POST /vehicles**
- **Description**: Adds a new vehicle to the system.
- **Method**: `POST`
- **URL**: `/vehicles`
- **Request Body**:
  ```json
  {
    "licensePlate": "ABC123",
    "model": "Toyota Corolla",
    "rentalRate": 30
  }
  ```
- **Response**:
  - **201 Created**: `"Vehicle added successfully with ID: {vehicleID}"`
  - **400 Bad Request**: `"Invalid input"`
  - **500 Internal Server Error**: `"Failed to insert vehicle"`

#### **3. POST /vehicle-status**
- **Description**: Adds a new vehicle status record.
- **Method**: `POST`
- **URL**: `/vehicle-status`
- **Request Body**:
  ```json
  {
    "vehicleID": 1,
    "status": "Available",
    "timestamp": "2024-12-08T14:00:00Z"
  }
  ```
- **Response**:
  - **201 Created**: `"Vehicle status added successfully"`
  - **400 Bad Request**: `"Invalid input"`
  - **500 Internal Server Error**: `"Failed to insert vehicle status"`

#### **4. POST /vehicle-booking**
- **Description**: Books a vehicle for a user.
- **Method**: `POST`
- **URL**: `/vehicle-booking`
- **Request Body**:
  ```json
  {
    "vehicleID": 1,
    "userID": 123,
    "startDate": "2024-12-10T10:00:00Z",
    "endDate": "2024-12-12T10:00:00Z",
    "totalPrice": 60
  }
  ```
- **Response**:
  - **201 Created**: `"Vehicle booked successfully"`
  - **400 Bad Request**: `"Invalid input"`
  - **500 Internal Server Error**: `"Failed to book vehicle"`

#### **5. POST /payment**
- **Description**: Processes a payment for a vehicle booking.
- **Method**: `POST`
- **URL**: `/payment`
- **Request Body**:
  ```json
  {
    "bookingID": 1,
    "paymentAmount": 60,
    "paymentStatus": "Successful"
  }
  ```
- **Response**:
  - **200 OK**: `"Payment processed successfully"`
  - **400 Bad Request**: `"Invalid payment data"`
  - **500 Internal Server Error**: `"Failed to process payment"`

#### **6. GET /vehicle/{vehicleID}**
- **Description**: Retrieves information about a specific vehicle.
- **Method**: `GET`
- **URL**: `/vehicle/{vehicleID}`
- **Response**:
  - **200 OK**: 
    ```json
    {
      "vehicleID": 1,
      "licensePlate": "ABC123",
      "model": "Toyota Corolla",
      "rentalRate": 30
    }
    ```
  - **404 Not Found**: `"Vehicle not found"`

#### **7. GET /vehicle-status/{vehicleID}**
- **Description**: Retrieves the status of a specific vehicle.
- **Method**: `GET`
- **URL**: `/vehicle-status/{vehicleID}`
- **Response**:
  - **200 OK**:
    ```json
    {
      "vehicleID": 1,
      "status": "Available",
      "timestamp": "2024-12-08T14:00:00Z"
    }
    ```
  - **404 Not Found**: `"Vehicle status not found"`

---

### **Error Codes**
- **400 Bad Request**: Invalid request input or parameters.
- **404 Not Found**: Resource not found (e.g., vehicle, booking, or status).
- **500 Internal Server Error**: An error occurred on the server (e.g., database issues).

---
