# API Routes Documentation

This document provides comprehensive documentation of all API routes, including authentication requirements, role/permission access, request/response formats, and example usage.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All protected routes require a valid JWT access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Role-Based Access Control (RBAC)

The API uses role-based access control. Default roles include:

| Role | Description |
|------|-------------|
| `admin` | Full system access |
| `manager` | Administrative access for business operations |
| `customer_experience` | Customer support access |
| `user` | Standard authenticated user |

---

## Response Format

### Success Response
```json
{
  "data": { ... },
  "meta": { ... }
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

### Pagination Metadata
```json
{
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

## Public Routes (No Authentication Required)

### Health Check

#### GET /health

Check if the API is running.

**Authentication:** None

**Response (200):**
```json
{
  "status": "ok"
}
```

---

## Authentication Routes

### POST /api/v1/auth/register

Register a new user account.

**Authentication:** None

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201):**
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
      "created_at": "2025-01-18T10:00:00Z",
      "updated_at": "2025-01-18T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "random-secure-token",
    "expires_at": "2025-01-18T10:15:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `409` - Email already exists

---

### POST /api/v1/auth/login

Authenticate with email and password.

**Authentication:** None

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response (200):**
```json
{
  "data": {
    "user": { ... },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "random-secure-token",
    "expires_at": "2025-01-18T10:15:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid credentials
- `403` - Account is inactive

---

### POST /api/v1/auth/refresh

Obtain a new access token using a refresh token.

**Authentication:** None

**Request Body:**
```json
{
  "refresh_token": "random-secure-token"
}
```

**Response (200):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "new-random-secure-token",
    "expires_at": "2025-01-18T10:15:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid refresh token

---

### GET /api/v1/auth/google

Get Google OAuth authorization URL.

**Authentication:** None

**Query Parameters:**
- `state` (optional) - CSRF protection state parameter

**Response (200):**
```json
{
  "data": {
    "url": "https://accounts.google.com/o/oauth2/auth?...",
    "state": "random-state-value"
  }
}
```

---

### GET /api/v1/auth/google/callback

Handle Google OAuth callback.

**Authentication:** None

**Query Parameters:**
- `code` (required) - Authorization code from Google
- `state` (required) - CSRF state parameter
- `error` (optional) - Error from Google if authorization failed

**Response (200):**
```json
{
  "data": {
    "user": { ... },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "random-secure-token",
    "expires_at": "2025-01-18T10:15:00Z"
  }
}
```

**Errors:**
- `400` - Missing authorization code or OAuth error

---

### GET /api/v1/auth/profile

Retrieve the authenticated user's profile.

**Authentication:** Required (any authenticated user)

**Permissions:** Any authenticated user

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
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

**Errors:**
- `401` - Authentication required
- `404` - User not found

---

### POST /api/v1/auth/logout

Revoke all refresh tokens for the current user.

**Authentication:** Required (any authenticated user)

**Permissions:** Any authenticated user

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (204):** No content

**Errors:**
- `401` - Authentication required

---

## Catalog Routes (Public)

### GET /api/v1/catalog/products

Retrieve a paginated list of products.

**Authentication:** None

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Products per page
- `keyword` (optional) - Search by product name or description

**Example:**
```
GET /api/v1/catalog/products?page=1&page_size=20&keyword=laptop
```

**Response (200):**
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
      "SalePrice": {
        "amount": 89999,
        "currency": "USD"
      },
      "status": "active",
      "brand_id": "brand-1",
      "category_id": "cat-1",
      "images": ["https://example.com/laptop.jpg"],
      "attributes": {"color": "silver", "ram": "16GB"},
      "created_at": "2025-01-18T10:00:00Z",
      "updated_at": "2025-01-18T10:00:00Z"
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

### GET /api/v1/catalog/products/:id

Retrieve details of a specific product.

**Authentication:** None

**Path Parameters:**
- `id` (required) - Product ID

**Example:**
```
GET /api/v1/catalog/products/prod-1
```

**Response (200):**
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
    "SalePrice": {
      "amount": 89999,
      "currency": "USD"
    },
    "status": "active",
    "brand_id": "brand-1",
    "category_id": "cat-1",
    "created_at": "2025-01-18T10:00:00Z",
    "updated_at": "2025-01-18T10:00:00Z"
  }
}
```

**Errors:**
- `400` - Product ID is required
- `404` - Product not found

---

### GET /api/v1/catalog/products/category/:id

Retrieve products in a specific category with pagination.

**Authentication:** None

**Path Parameters:**
- `id` (required) - Category ID

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:**
```
GET /api/v1/catalog/products/category/cat-1?page=1&page_size=10
```

**Response (200):**
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

### GET /api/v1/catalog/categories

Retrieve all categories with pagination.

**Authentication:** None

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response (200):**
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
      "created_at": "2025-01-18T10:00:00Z",
      "updated_at": "2025-01-18T10:00:00Z"
    }
  ],
  "meta": { ... }
}
```

---

### GET /api/v1/catalog/brands

Retrieve all brands with pagination.

**Authentication:** None

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response (200):**
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
      "created_at": "2025-01-18T10:00:00Z",
      "updated_at": "2025-01-18T10:00:00Z"
    }
  ],
  "meta": { ... }
}
```

---

## Cart Routes (Protected - Any Authenticated User)

All cart routes require authentication. Users can only access their own cart.

### GET /api/v1/cart

Retrieve the current user's cart.

**Authentication:** Required

**Permissions:** Any authenticated user (owns their cart)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
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
    "created_at": "2025-01-18T10:00:00Z",
    "updated_at": "2025-01-18T10:05:00Z"
  }
}
```

**Errors:**
- `401` - Authentication required

---

### POST /api/v1/cart/items

Add a product to the cart.

**Authentication:** Required

**Permissions:** Any authenticated user

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

**Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Errors:**
- `400` - Invalid request body or product out of stock
- `401` - Authentication required

---

### PATCH /api/v1/cart/items/:id

Update the quantity of an item in the cart.

**Authentication:** Required

**Permissions:** Any authenticated user (owns their cart)

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

**Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Errors:**
- `400` - Invalid request body or item ID required
- `401` - Authentication required
- `404` - Item not found in cart

---

### DELETE /api/v1/cart/items/:id

Remove an item from the cart.

**Authentication:** Required

**Permissions:** Any authenticated user (owns their cart)

**Path Parameters:**
- `id` (required) - Cart item ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    /* Updated cart object */
  }
}
```

**Errors:**
- `400` - Item ID is required
- `401` - Authentication required
- `404` - Item not found in cart

---

### DELETE /api/v1/cart

Remove all items from the cart.

**Authentication:** Required

**Permissions:** Any authenticated user (owns their cart)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    /* Empty cart object */
  }
}
```

**Errors:**
- `401` - Authentication required

---

## Order Routes (Protected)

### POST /api/v1/orders

Create an order from the current user's cart.

**Authentication:** Required

**Permissions:** Any authenticated user

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

**Response (201):**
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
    "created_at": "2025-01-18T10:00:00Z",
    "updated_at": "2025-01-18T10:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body, cart is empty, or invalid address
- `401` - Authentication required

---

### GET /api/v1/orders

Retrieve the current user's orders with pagination.

**Authentication:** Required

**Permissions:** Any authenticated user (views their own orders)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:**
```
GET /api/v1/orders?page=1&page_size=10
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": [
    {
      /* Order object */
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total_items": 5,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  }
}
```

**Errors:**
- `401` - Authentication required

---

### GET /api/v1/orders/:id

Retrieve details of a specific order.

**Authentication:** Required

**Permissions:**
- Order owner (user_id matches authenticated user)
- **OR** Users with role: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Order ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    /* Order object */
  }
}
```

**Errors:**
- `400` - Order ID is required
- `401` - Authentication required
- `403` - You don't have permission to view this order
- `404` - Order not found

---

## Admin Routes

All admin routes require authentication AND one of the following roles:
- `admin`
- `manager`
- `customer_experience`

---

## Role Management

### GET /api/v1/admin/roles

List all roles.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "roles": [
      {
        "id": "role-id",
        "name": "admin",
        "description": "Full system access"
      }
    ]
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

### POST /api/v1/admin/roles

Create a new role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "inventory_manager",
  "description": "Manages inventory levels"
}
```

**Response (201):**
```json
{
  "data": {
    "role": {
      "id": "role-id",
      "name": "inventory_manager",
      "description": "Manages inventory levels"
    }
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Authentication required
- `403` - Insufficient permissions

---

### GET /api/v1/admin/roles/:id

Get a specific role by ID.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "role": {
      "id": "role-id",
      "name": "admin",
      "description": "Full system access"
    }
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Role not found

---

### PUT /api/v1/admin/roles/:id

Update an existing role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "inventory_admin",
  "description": "Full inventory access"
}
```

**Response (200):**
```json
{
  "data": {
    "role": {
      "id": "role-id",
      "name": "inventory_admin",
      "description": "Full inventory access"
    }
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Role not found

---

### DELETE /api/v1/admin/roles/:id

Delete a role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (204):** No content

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

### GET /api/v1/admin/roles/:id/permissions

Get all permissions granted to a role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "permissions": [
      {
        "id": "perm-id",
        "name": "manage_products",
        "resource": "products",
        "action": "manage",
        "description": "Full access to products"
      }
    ]
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

### POST /api/v1/admin/roles/:id/permissions

Grant a permission to a role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "permission_id": "perm-id"
}
```

**Response (200):**
```json
{
  "data": {
    "message": "Permission granted successfully"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Authentication required
- `403` - Insufficient permissions

---

### DELETE /api/v1/admin/roles/:id/permissions/:permId

Revoke a permission from a role.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Role ID
- `permId` (required) - Permission ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (204):** No content

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

## Permission Management

### GET /api/v1/admin/permissions

List all permissions.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "permissions": [
      {
        "id": "perm-id",
        "name": "manage_products",
        "resource": "products",
        "action": "manage",
        "description": "Full access to products"
      }
    ]
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

### POST /api/v1/admin/permissions

Create a new permission.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "manage_inventory",
  "resource": "inventory",
  "action": "manage",
  "description": "Full access to inventory management"
}
```

**Response (201):**
```json
{
  "data": {
    "permission": {
      "id": "perm-id",
      "name": "manage_inventory",
      "resource": "inventory",
      "action": "manage",
      "description": "Full access to inventory management"
    }
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Authentication required
- `403` - Insufficient permissions

---

### GET /api/v1/admin/permissions/:id

Get a specific permission by ID.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Permission ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "permission": {
      "id": "perm-id",
      "name": "manage_products",
      "resource": "products",
      "action": "manage",
      "description": "Full access to products"
    }
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Permission not found

---

### PUT /api/v1/admin/permissions/:id

Update an existing permission.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Permission ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "manage_inventory",
  "resource": "inventory",
  "action": "manage",
  "description": "Updated description"
}
```

**Response (200):**
```json
{
  "data": {
    "permission": {
      "id": "perm-id",
      "name": "manage_inventory",
      "resource": "inventory",
      "action": "manage",
      "description": "Updated description"
    }
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Permission not found

---

### DELETE /api/v1/admin/permissions/:id

Delete a permission.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - Permission ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (204):** No content

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

## User Role Assignments

### GET /api/v1/admin/users/:id/roles

Get all roles assigned to a user.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - User ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "data": {
    "roles": ["admin", "user"]
  }
}
```

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions

---

### POST /api/v1/admin/users/:id/roles

Assign a role to a user.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - User ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "role_id": "role-id"
}
```
or
```json
{
  "role_name": "admin"
}
```

**Response (200):**
```json
{
  "data": {
    "message": "Role assigned successfully"
  }
}
```

**Errors:**
- `400` - Invalid request body or either role_id or role_name is required
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Role not found

---

### DELETE /api/v1/admin/users/:id/roles/:roleId

Remove a role from a user.

**Authentication:** Required

**Permissions:** Role required: `admin`, `manager`, or `customer_experience`

**Path Parameters:**
- `id` (required) - User ID
- `roleId` (required) - Role ID

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (204):** No content

**Errors:**
- `401` - Authentication required
- `403` - Insufficient permissions
- `404` - Role not found

---

## Route Summary Table

| Method | Path | Auth | Roles/Permissions |
|--------|------|------|-------------------|
| GET | /health | No | - |
| POST | /api/v1/auth/register | No | - |
| POST | /api/v1/auth/login | No | - |
| POST | /api/v1/auth/refresh | No | - |
| GET | /api/v1/auth/google | No | - |
| GET | /api/v1/auth/google/callback | No | - |
| GET | /api/v1/auth/profile | Yes | Any authenticated user |
| POST | /api/v1/auth/logout | Yes | Any authenticated user |
| GET | /api/v1/catalog/products | No | - |
| GET | /api/v1/catalog/products/:id | No | - |
| GET | /api/v1/catalog/products/category/:id | No | - |
| GET | /api/v1/catalog/categories | No | - |
| GET | /api/v1/catalog/brands | No | - |
| GET | /api/v1/cart | Yes | Any authenticated user |
| POST | /api/v1/cart/items | Yes | Any authenticated user |
| PATCH | /api/v1/cart/items/:id | Yes | Any authenticated user |
| DELETE | /api/v1/cart/items/:id | Yes | Any authenticated user |
| DELETE | /api/v1/cart | Yes | Any authenticated user |
| POST | /api/v1/orders | Yes | Any authenticated user |
| GET | /api/v1/orders | Yes | Any authenticated user |
| GET | /api/v1/orders/:id | Yes | Owner OR admin/manager/customer_experience |
| GET | /api/v1/admin/roles | Yes | admin, manager, customer_experience |
| POST | /api/v1/admin/roles | Yes | admin, manager, customer_experience |
| GET | /api/v1/admin/roles/:id | Yes | admin, manager, customer_experience |
| PUT | /api/v1/admin/roles/:id | Yes | admin, manager, customer_experience |
| DELETE | /api/v1/admin/roles/:id | Yes | admin, manager, customer_experience |
| GET | /api/v1/admin/roles/:id/permissions | Yes | admin, manager, customer_experience |
| POST | /api/v1/admin/roles/:id/permissions | Yes | admin, manager, customer_experience |
| DELETE | /api/v1/admin/roles/:id/permissions/:permId | Yes | admin, manager, customer_experience |
| GET | /api/v1/admin/permissions | Yes | admin, manager, customer_experience |
| POST | /api/v1/admin/permissions | Yes | admin, manager, customer_experience |
| GET | /api/v1/admin/permissions/:id | Yes | admin, manager, customer_experience |
| PUT | /api/v1/admin/permissions/:id | Yes | admin, manager, customer_experience |
| DELETE | /api/v1/admin/permissions/:id | Yes | admin, manager, customer_experience |
| GET | /api/v1/admin/users/:id/roles | Yes | admin, manager, customer_experience |
| POST | /api/v1/admin/users/:id/roles | Yes | admin, manager, customer_experience |
| DELETE | /api/v1/admin/users/:id/roles/:roleId | Yes | admin, manager, customer_experience |

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | OK - Request succeeded |
| `201` | Created - Resource created successfully |
| `204` | No Content - Request succeeded with no response body |
| `400` | Bad Request - Invalid request data |
| `401` | Unauthorized - Authentication required or invalid token |
| `403` | Forbidden - Insufficient permissions |
| `404` | Not Found - Resource not found |
| `409` | Conflict - Resource already exists |
| `500` | Internal Server Error - Server error |
