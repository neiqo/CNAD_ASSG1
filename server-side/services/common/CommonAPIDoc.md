# Common Service API Documentation

## Overview
The Common Service API provides endpoints to manage member benefits, promotions, and service status. It connects to a MySQL database to fetch and manage member benefits and promotions.

### Base URL
The service is available at `http://localhost:5003`.

## Endpoints

### 1. **GET /api/v1/status**
Check the connection status of the Common Service to the database.

#### Response
- **200 OK**: Common Service is connected to the database.
- **500 Internal Server Error**: Database connection failed.

#### Example Request:
```
GET /api/v1/status
```

#### Example Response:
```
"Common Service connected to the database successfully!"
```

---

### 2. **GET /api/v1/member-benefits/{membershipTier}**
Fetch member benefits for a specific membership tier.

#### Path Parameters
- `membershipTier` (string): The membership tier to fetch benefits for.

#### Response
- **200 OK**: Member benefits retrieved successfully.
- **404 Not Found**: No benefits found for the specified membership tier.
- **500 Internal Server Error**: Error fetching member benefits.

#### Example Request:
```
GET /api/v1/member-benefits/Gold
```

#### Example Response:
```
[
  {
    "BenefitID": 1,
    "Name": "Free Shipping",
    "Description": "Get free shipping on all orders."
  },
  {
    "BenefitID": 2,
    "Name": "Exclusive Discounts",
    "Description": "Access to exclusive discounts on products."
  }
]
```

---

### 3. **GET /api/v1/promotions**
Retrieve all available promotions.

#### Response
- **200 OK**: Promotions retrieved successfully.
- **404 Not Found**: No promotions found.
- **500 Internal Server Error**: Error fetching promotions.

#### Example Request:
```
GET /api/v1/promotions
```

#### Example Response:
```
[
  {
    "PromotionID": 1,
    "Name": "Black Friday Sale",
    "Description": "Get 50% off on all electronics.",
    "Discount": 50,
    "IfPercentage": true
  },
  {
    "PromotionID": 2,
    "Name": "Winter Discount",
    "Description": "Get $10 off on all orders above $50.",
    "Discount": 10,
    "IfPercentage": false
  }
]
```
