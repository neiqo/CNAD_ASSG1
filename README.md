
# Vehicle Rental Services API 

This repository contains a set of microservices for a vehicle rental system stored under the server-side folder and client-side access stored on the client-side folder, which includes the following services:

- **User Service**: Manages user registrations, authentication, and profiles.
- **Vehicle Service**: Handles vehicle details, availability, and booking management.
- **Common Service**: Provides shared utilities and common functionalities across other services.
- **Billing Service**: Handles the calculation of charges, payment processing, and invoice generation.

## Services Overview 

![enter image description here](https://i.imgur.com/mPZZw7b.png)

### 1. **User Service** (Port 5001)
The User Service is responsible for managing user accounts, including registration, authentication, and profile management.

#### Key Features:
- User encrypted registration and login.
- Profile management.
- Retrieves user details based on their email
- Sends verification code to email

### 2. **Vehicle Service** (Port 5002)
The Vehicle Service manages the vehicle details, status, and booking processes.

#### Key Features:
- Handles fetching and managing car and booking data.
- Provides APIs for retrieving vehicle information, searching for available vehicles, and managing bookings.
- Interacts with the payment-service and user-service for booking and billing
- Provides booking history for users

### 3. **Common Service** (Port 5003)
The Common Service provides shared utilities and components that are used across all other services. This includes common data formats, validation logic, and utilities.

#### Key Features:
- Shared data so billing and users can use the data without slowing either one of the services down due to too many requests

### 4. **Billing Service** (5004)
The Billing Service is responsible for managing all aspects related to vehicle bookings and payments.

#### Key Features:
- Payment processing and calculation.
- Invoice generation and sending to the users' email.

### 5. Independent Services and Databases
Each microservice operates independently, with its own isolated codebase and database. This architecture promotes a modular, scalable, and maintainable system where each service is responsible for a specific set of tasks.
---

## Getting Started

Follow the instructions below to set up and run the microservices locally.

### Prerequisites

- 
- MySQL
- Postman or similar tool for testing API endpoints

### Setting Up the Services


1. **Run the services:**

  To run each service locally you can run the RunAllServices.ps1 powershell script

2. **Access the services:**
   - The services will run on different ports. By default, they are set to:
     - User Service: `http://localhost:5001`
     - Vehicle Service: `http://localhost:5002`
     - Common Service: `http://localhost:5003`
     - Billing Service: `http://localhost:5004`

3. **Testing the services:**
   You can test the API endpoints using Postman or curl to send requests to the services.

---

## Detailed API Documentation

For each service, refer to their individual documentation for details on API endpoints, request/response formats, and error codes that is stored in its folder e.g: server-side/services/user/UserAPIDoc.md.

---
