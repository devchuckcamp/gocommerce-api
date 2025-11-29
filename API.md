# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Pagination

All list endpoints support pagination using query parameters:

### Pagination Parameters
- `page` - Page number (default: 1, min: 1)
- `page_size` - Items per page (default: 20, min: 1, max: 100)

### Pagination Response
```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "page_size": 20,
    "total_items": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

## Response Format

All responses follow a consistent structure:

### Success Response
```json
{
  "data": { ... },
  "meta": { ... }  // Pagination metadata for list endpoints
}
```

### Error Response
```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message"
  }
}
```

## HTTP Status Codes

- `200` - OK: Request succeeded
- `201` - Created: Resource created successfully
- `204` - No Content: Request succeeded with no response body
- `400` - Bad Request: Invalid request data
- `401` - Unauthorized: Authentication required or invalid token
- `403` - Forbidden: Insufficient permissions
- `404` - Not Found: Resource not found
- `409` - Conflict: Resource already exists
- `500` - Internal Server Error: Server error

---

## Authentication Endpoints

### Register User
Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Success Response (201):**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "active": true,
      "email_verified": false,
      "created_at": "2025-11-29T10:00:00Z",
      "updated_at": "2025-11-29T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "random-secure-token",
    "expires_at": "2025-11-29T10:15:00Z"
  }
}
```

**Error Responses:**
- `400` - Invalid request body
- `409` - Email already exists

---

### Login
Authenticate with email and password.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Success Response (200):**
```json
{
  "data": {
    "user": { ... },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "random-secure-token",
    "expires_at": "2025-11-29T10:15:00Z"
  }
}
```

**Error Responses:**
- `400` - Invalid request body
- `401` - Invalid credentials
- `403` - Account is inactive

---

### Get Profile
Retrieve the authenticated user's profile.

**Endpoint:** `GET /auth/profile`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "active": true,
      "email_verified": false
    },
    "roles": ["user"]
  }
}
```

**Error Responses:**
- `401` - Authentication required
- `404` - User not found

---

### Refresh Token
Obtain a new access token using a refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "random-secure-token"
}
```

**Success Response (200):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "new-random-secure-token",
    "expires_at": "2025-11-29T10:15:00Z"
  }
}
```

**Error Responses:**
- `400` - Invalid request body
- `401` - Invalid refresh token

---

### Logout
Revoke all refresh tokens for the current user.

**Endpoint:** `POST /auth/logout`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (204):**
No content

**Error Responses:**
- `401` - Authentication required

---

## Catalog Endpoints (Public)

### List Products
Retrieve a paginated list of products.

**Endpoint:** `GET /catalog/products`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Products per page
- `keyword` (optional) - Search by product name or description

**Examples:**
- `GET /catalog/products?page=1&page_size=20`
- `GET /catalog/products?keyword=laptop&page=1`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "prod-1",
      "sku": "LAPTOP-001",
      "name": "Professional Laptop",
      "description": "High-performance laptop for professionals",
      "base_price": {
        "amount": 99999,
        "currency": "USD"
      },
      "status": "active",
      "brand_id": "brand-1",
      "category_id": "cat-1",
      "images": ["https://example.com/laptop.jpg"],
      "attributes": {"color": "silver", "ram": "16GB"},
      "created_at": "2025-11-29T10:00:00Z",
      "updated_at": "2025-11-29T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 20,
    "total_items": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

---

### Get Product
Retrieve details of a specific product.

**Endpoint:** `GET /catalog/products/:id`

**Path Parameters:**
- `id` (required) - Product ID

**Example:** `GET /catalog/products/prod-1`

**Success Response (200):**
```json
{
  "data": {
    "id": "prod-1",
    "sku": "LAPTOP-001",
    "name": "Professional Laptop",
    "description": "High-performance laptop for professionals",
    "base_price": {
      "amount": 99999,
      "currency": "USD"
    },
    "status": "active",
    "brand_id": "brand-1",
    "category_id": "cat-1",
    "image_url": "",
    "metadata": {},
    "created_at": "2025-11-29T10:00:00Z",
    "updated_at": "2025-11-29T10:00:00Z"
  }
}
```

**Error Responses:**
- `400` - Product ID is required
- `404` - Product not found

---

### Get Products by Category
Retrieve products in a specific category with pagination.

**Endpoint:** `GET /catalog/products/category/:id`

**Path Parameters:**
- `id` (required) - Category ID

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:** `GET /catalog/products/category/cat-1?page=1&page_size=10`

**Success Response (200):**
```json
{
  "data": [
    { /* product object */ }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total_items": 45,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

---

### List Categories
Retrieve all categories.

**Endpoint:** `GET /catalog/categories`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "cat-1",
      "name": "Electronics",
      "slug": "electronics",
      "description": "Electronic devices and gadgets",
      "parent_id": null,
      "image_url": "",
      "active": true,
      "created_at": "2025-11-29T10:00:00Z",
      "updated_at": "2025-11-29T10:00:00Z"
    }
  ]
}
```

---

### List Brands
Retrieve all brands.

**Endpoint:** `GET /catalog/brands`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "brand-1",
      "name": "TechCorp",
      "slug": "techcorp",
      "description": "Leading technology manufacturer",
      "logo_url": "",
      "active": true,
      "created_at": "2025-11-29T10:00:00Z",
      "updated_at": "2025-11-29T10:00:00Z"
    }
  ]
}
```

---

## Cart Endpoints (Protected)

All cart endpoints require authentication.

### Get Cart
Retrieve the current user's cart.

**Endpoint:** `GET /cart`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": {
    "id": "cart-id",
    "user_id": "user-id",
    "session_id": "",
    "items": [
      {
        "id": "item-id",
        "product_id": "prod-1",
        "variant_id": null,
        "sku": "LAPTOP-001",
        "name": "Professional Laptop",
        "price": {
          "amount": 99999,
          "currency": "USD"
        },
        "quantity": 1,
        "attributes": {}
      }
    ],
    "created_at": "2025-11-29T10:00:00Z",
    "updated_at": "2025-11-29T10:05:00Z"
  }
}
```

**Error Responses:**
- `401` - Authentication required

---

### Add Item to Cart
Add a product to the cart.

**Endpoint:** `POST /cart/items`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "product_id": "prod-1",
  "variant_id": null,
  "quantity": 2,
  "attributes": {
    "color": "Black",
    "size": "15-inch"
  }
}
```

**Success Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Error Responses:**
- `400` - Invalid request body or product out of stock
- `401` - Authentication required

---

### Update Item Quantity
Update the quantity of an item in the cart.

**Endpoint:** `PATCH /cart/items/:id`

**Authentication:** Required

**Path Parameters:**
- `id` (required) - Cart item ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "quantity": 3
}
```

**Success Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Error Responses:**
- `400` - Invalid request body or item ID required
- `401` - Authentication required
- `404` - Item not found in cart

---

### Remove Item from Cart
Remove an item from the cart.

**Endpoint:** `DELETE /cart/items/:id`

**Authentication:** Required

**Path Parameters:**
- `id` (required) - Cart item ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Error Responses:**
- `400` - Item ID is required
- `401` - Authentication required
- `404` - Item not found in cart

---

### Clear Cart
Remove all items from the cart.

**Endpoint:** `DELETE /cart`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": {
    /* Empty cart object */
  }
}
```

---

## Order Endpoints (Protected)

All order endpoints require authentication.

### Create Order
Create an order from the current user's cart.

**Endpoint:** `POST /orders`

**Authentication:** Required

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "shipping_address": {
    "first_name": "John",
    "last_name": "Doe",
    "company": "Tech Corp",
    "address1": "123 Main St",
    "address2": "Apt 4B",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "US",
    "phone_number": "555-0100"
  },
  "billing_address": {
    /* Same format as shipping_address, optional */
  },
  "payment_method_id": "pm_123",
  "promotion_codes": ["SAVE10"],
  "shipping_method_id": "ship_standard",
  "notes": "Please deliver after 5 PM"
}
```

**Success Response (201):**
```json
{
  "data": {
    "id": "order-id",
    "order_number": "ORD-12345678",
    "user_id": "user-id",
    "status": "pending",
    "items": [
      {
        "id": "item-id",
        "product_id": "prod-1",
        "sku": "LAPTOP-001",
        "name": "Professional Laptop",
        "unit_price": {
          "amount": 99999,
          "currency": "USD"
        },
        "quantity": 1,
        "discount_amount": {
          "amount": 0,
          "currency": "USD"
        },
        "tax_amount": {
          "amount": 8750,
          "currency": "USD"
        },
        "total": {
          "amount": 108749,
          "currency": "USD"
        }
      }
    ],
    "shipping_address": { /* address object */ },
    "billing_address": { /* address object */ },
    "subtotal": {
      "amount": 99999,
      "currency": "USD"
    },
    "discount_total": {
      "amount": 0,
      "currency": "USD"
    },
    "tax_total": {
      "amount": 8750,
      "currency": "USD"
    },
    "shipping_total": {
      "amount": 0,
      "currency": "USD"
    },
    "total": {
      "amount": 108749,
      "currency": "USD"
    },
    "notes": "Please deliver after 5 PM",
    "created_at": "2025-11-29T10:00:00Z",
    "updated_at": "2025-11-29T10:00:00Z"
  }
}
```

**Error Responses:**
- `400` - Invalid request body, cart is empty, or invalid address
- `401` - Authentication required

---

### List Orders
Retrieve the current user's orders with pagination.

**Endpoint:** `GET /orders`

**Authentication:** Required

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:** `GET /orders?page=1&page_size=10`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": [
    {
      /* Order object */
    }
  ]
}
```

---

### Get Order
Retrieve details of a specific order.

**Endpoint:** `GET /orders/:id`

**Authentication:** Required

**Path Parameters:**
- `id` (required) - Order ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "data": {
    /* Order object */
  }
}
```

**Error Responses:**
- `400` - Order ID is required
- `401` - Authentication required
- `403` - You don't have permission to view this order
- `404` - Order not found

---

## Health Check

### Health Check
Check if the API is running.

**Endpoint:** `GET /health`

**Success Response (200):**
```json
{
  "status": "ok"
}
```
