# E-Commerce API with Gin, GORM, goauthx, and gocommerce

A production-ready Go E-Commerce API that integrates:
- **Gin** - Fast HTTP web framework
- **GORM** - Feature-rich ORM for database operations
- **[goauthx](https://github.com/devchuckcamp/goauthx)** - Complete authentication & authorization solution
- **[gocommerce](https://github.com/devchuckcamp/gocommerce)** - E-commerce domain logic library

## Core Packages

### 🔐 [goauthx v0.0.3](https://github.com/devchuckcamp/goauthx)
Authentication and authorization library with:
- JWT token management (access & refresh)
- Google OAuth 2.0 integration
- Role-based access control (RBAC) with permissions
- Admin API for role/permission management
- User management with email verification
- Password reset functionality
- Multi-database support (PostgreSQL, MySQL, SQL Server)

### 🛒 [gocommerce v0.0.5](https://github.com/devchuckcamp/gocommerce)
E-commerce domain logic library with:
- Catalog management (products, variants, categories, brands)
- Shopping cart with persistence
- Order processing and management
- Pricing and promotions with date-windowed sale prices
- Tax calculations
- **Inventory management** (stock levels, reservations, reorder points)
- **Supplier management** (multi-supplier per product, cost tracking)
- **Inventory audit logging** (stock_in, stock_out, adjustments, transfers)
- Clean domain-driven design

## Features

### Authentication & Authorization (powered by goauthx)
- ✅ User registration and login with JWT tokens
- ✅ **Google OAuth 2.0** - Sign in with Google
- ✅ Access token & refresh token management
- ✅ Role-based access control (RBAC) with permissions
- ✅ **Admin API**: Role/permission management, user role assignments
- ✅ Protected routes with middleware
- ✅ User profile management

### E-Commerce (powered by gocommerce)
- ✅ **Catalog**: Products, variants, categories, and brands
- ✅ **Product Search**: Keyword search by name/description
- ✅ **Pagination**: All listing endpoints (products, categories, brands, orders)
- ✅ **Shopping Cart**: Add/update/remove items, cart persistence
- ✅ **Orders**: Create orders from cart, order history with pagination
- ✅ **Pricing**: Tax calculation, promotion support, date-windowed sale prices
- ✅ **Inventory Ready**: Database tables for stock levels, reservations, suppliers
- ✅ Clean domain-driven architecture

### Technical Features
- ✅ Multi-database support (PostgreSQL, MySQL, SQL Server)
- ✅ **Docker Deployment**: Multi-stage builds, health checks
- ✅ RESTful API design with consistent responses
- ✅ Pagination with metadata (page, total_items, has_next/prev)
- ✅ Structured logging and error handling
- ✅ CORS support
- ✅ Graceful shutdown
- ✅ Environment-based configuration

## Project Structure

```
github.com/devchuckcamp/gocommerce-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── database/
│   │   ├── database.go             # GORM connection setup
│   │   ├── models.go               # Database models
│   │   └── migrations.go           # Auto-migration & seeding
│   ├── repository/
│   │   ├── catalog.go              # Product/Category/Brand repositories
│   │   ├── cart.go                 # Cart repository
│   │   ├── orders.go               # Orders repository
│   │   └── pricing.go              # Promotion repository
│   ├── services/
│   │   ├── catalog.go              # Catalog service with search
│   │   ├── cart.go                 # Cart service (gocommerce wrapper)
│   │   ├── orders.go               # Order service (gocommerce wrapper)
│   │   ├── pricing.go              # Pricing service (gocommerce wrapper)
│   │   └── tax.go                  # Tax calculator implementation
│   ├── http/
│   │   ├── server.go               # HTTP server & route setup
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWT auth middleware
│   │   │   ├── logger.go           # Request logging
│   │   │   ├── recovery.go         # Panic recovery
│   │   │   └── cors.go             # CORS middleware
│   │   ├── handlers/
│   │   │   ├── admin.go            # Admin RBAC handlers (roles, permissions, users)
│   │   │   ├── auth.go             # Auth + Google OAuth handlers
│   │   │   ├── catalog.go          # Catalog handlers with pagination
│   │   │   ├── cart.go             # Cart handlers
│   │   │   └── orders.go           # Order handlers with pagination
│   │   └── response/
│   │       └── response.go         # API responses + pagination
│   └── utils/
│       └── id.go                   # ID generation utilities
├── .dockerignore                   # Docker build exclusions
├── .env.example                    # Environment template
├── .gitignore                      # Git exclusions
├── API.md                          # Complete API documentation
├── docker-compose.yml              # Docker orchestration
├── docker-start.sh                 # Docker deployment script
├── docker-stop.sh                  # Docker stop script
├── Dockerfile                      # Multi-stage build config
├── DOCKER.md                       # Docker deployment guide
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── GOOGLE_OAUTH.md                 # Google OAuth setup guide
├── README.md                       # This file
├── ROUTES.md                       # Complete API routes with permissions
├── setup.sh                        # Local development setup
├── stop.sh                         # Local process stop script
└── test-oauth.html                 # OAuth testing page
```

## Quick Start

### Prerequisites

- **Docker & Docker Compose** (recommended) OR
- Go 1.23+ and PostgreSQL/MySQL/SQL Server
- Git

### Option 1: Docker Deployment (Recommended)

```bash
# Clone the repository
git clone https://github.com/devchuckcamp/gocommerce-api.git
cd gocommerce-api

# Copy and configure environment
cp .env.example .env
# Edit .env with your Google OAuth credentials (optional)

# Start services with Docker
./docker-start.sh

# API will be available at http://localhost:8080

# Stop services
./docker-stop.sh
```

### Option 2: Local Development

```bash
# Clone the repository
git clone https://github.com/devchuckcamp/gocommerce-api.git
cd gocommerce-api

# Run setup script
./setup.sh

# Edit .env with your configuration
# IMPORTANT: Set JWT_SECRET (min 32 chars) and database credentials
nano .env

# Create database (PostgreSQL example)
psql -U postgres -c "CREATE DATABASE commerce;"

# Run the application (migrations run automatically)
go run cmd/api/main.go

# Stop with Ctrl+C or use
./stop.sh
```

3. **Setup environment variables**
```bash
cp .env.example .env
# Edit .env with your database credentials and JWT secret
```

4. **Generate a secure JWT secret**
```bash
# Use a random string generator or:
openssl rand -base64 32
```

5. **Set up your database**

Create a database for the application:

**PostgreSQL:**
```sql
CREATE DATABASE gocommerce;
```

**MySQL:**
```sql
CREATE DATABASE gocommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

6. **Run migrations and seed (optional)**

Migrations run automatically on startup using the official migration systems from both libraries:

- **goauthx migrations**: Creates authentication tables (users, roles, permissions, refresh_tokens, etc.)
- **gocommerce migrations**: Creates e-commerce tables (products, categories, brands, carts, orders, etc.)

To seed with sample data:
```bash
export SEED_DB=true
# or in .env file: SEED_DB=true
```

The seeding process will populate:
- Sample products (laptops, phones, tablets)
- Categories (electronics, computers, accessories)
- Brands (Apple, Dell, Lenovo, HP, Samsung)

7. **Run the application**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

> **Complete API documentation with request/response examples, authentication requirements, and role-based permissions is available in [ROUTES.md](ROUTES.md).**

### Endpoint Summary

| Category | Routes | Auth Required | Roles |
|----------|--------|---------------|-------|
| Health | `GET /health` | No | - |
| Auth | 6 endpoints | Mixed | - |
| Catalog | 5 endpoints | No | - |
| Cart | 5 endpoints | Yes | Any user |
| Orders | 3 endpoints | Yes | Any user / Admin |
| Admin RBAC | 16 endpoints | Yes | admin, manager, customer_experience |

**Total: 37 API endpoints**

For detailed documentation including:
- Request/response formats with examples
- Role and permission requirements per endpoint
- Error codes and handling

See **[ROUTES.md](ROUTES.md)**

## Usage Examples

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!"
  }'
```

Save the `access_token` from the response.

### 3. List Products

```bash
curl http://localhost:8080/api/v1/catalog/products
```

### 4. Add Item to Cart

```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "product_id": "prod-1",
    "quantity": 2
  }'
```

### 5. Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "shipping_address": {
      "first_name": "John",
      "last_name": "Doe",
      "address1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "US",
      "phone_number": "555-0100"
    }
  }'
```

## Database Support

The application supports multiple database engines through GORM:

### PostgreSQL (Recommended)
```env
DB_DRIVER=postgres
DB_DSN=postgres://user:password@localhost:5432/gocommerce?sslmode=disable
```

### MySQL
```env
DB_DRIVER=mysql
DB_DSN=user:password@tcp(localhost:3306)/gocommerce?parseTime=true
```

### SQL Server
```env
DB_DRIVER=sqlserver
DB_DSN=sqlserver://user:password@localhost:1433?database=gocommerce
```

## Architecture Highlights

### Clean Architecture
- **Domain Layer**: Business logic from `gocommerce` library
- **Repository Layer**: GORM implementations of repository interfaces
- **Service Layer**: Orchestration of domain services
- **HTTP Layer**: Gin handlers and middleware

### Dependency Injection
All dependencies are injected at startup, making the application:
- Easy to test
- Easy to swap implementations
- Clear dependency graph

### Error Handling
Consistent error responses across all endpoints:
```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message"
  }
}
```

Success responses:
```json
{
  "data": { ... },
  "meta": { ... }
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

### Docker Support

**Full Docker deployment is available!** See [DOCKER.md](DOCKER.md) for complete guide.

#### Quick Docker Start

Run both database and API in containers:
```bash
docker-compose up --build
```

This starts:
- **PostgreSQL database** (postgres:16-alpine)
- **Go API service** (built from Dockerfile)
- Automatic migrations and health checks
- Access API at http://localhost:8080

#### Database Only

Run just PostgreSQL database:
```bash
docker-compose up postgres -d
```

Configuration:
- Username: `commerce`
- Password: `commerce`
- Database: `commerce`
- Port: `5432`

For detailed Docker documentation, see [DOCKER.md](DOCKER.md).

## Migration System

This project uses the official migration capabilities from both **goauthx** and **gocommerce** libraries:

### goauthx Migrations
Creates authentication and authorization tables:
- `users` - User accounts
- `roles` - User roles
- `permissions` - Permission definitions
- `user_roles` - User-role assignments
- `role_permissions` - Role-permission assignments
- `refresh_tokens` - JWT refresh tokens
- `email_verifications` - Email verification tokens
- `password_resets` - Password reset tokens
- `schema_migrations` - Migration tracking

### gocommerce Migrations
Creates e-commerce domain tables:
- `products` - Product catalog
- `categories` - Product categories
- `brands` - Product brands
- `carts` - Shopping carts
- `cart_items` - Cart line items
- `orders` - Customer orders
- `product_prices` - Date-windowed sale prices
- `suppliers` - Supplier management (v0.0.5)
- `product_suppliers` - Product-supplier relationships (v0.0.5)
- `inventory_levels` - Stock tracking (on_hand, reserved, available) (v0.0.5)
- `inventory_suppliers` - Supplier-specific inventory (v0.0.5)
- `inventory_activities` - Inventory audit logging (v0.0.5)
- `gocommerce_migrations` - Migration tracking

### How Migrations Work

1. **Automatic Execution**: Migrations run automatically when the application starts
2. **Idempotent**: Safe to run multiple times - only applies new migrations
3. **Version Tracking**: Both systems track which migrations have been applied
4. **Transaction Safety**: Each migration runs in a transaction

### Running Migrations

**Migrations run automatically** - just start the application:

```bash
# Run migrations only (without seeding)
go run cmd/api/main.go
# or
./bin/api.exe
```

**To run migrations AND seed sample data:**

```bash
# Set environment variable
export SEED_DB=true

# Then run the application
go run cmd/api/main.go
```

Or update your `.env` file:
```env
SEED_DB=true
```

Then start the application:
```bash
go run cmd/api/main.go
```

The startup logs will show:
```
Running goauthx migrations...
✓ goauthx migrations completed successfully
Running gocommerce migrations...
✓ gocommerce migrations completed successfully
Seeding database with sample data...
✓ Database seeded successfully
```

### Manual Migration Commands

If you need to check migration status or troubleshoot:

```bash
# View all tables created
docker exec goshop-postgres psql -U commerce -d commerce -c "\dt"

# Check migration history
docker exec goshop-postgres psql -U commerce -d commerce -c "SELECT * FROM schema_migrations;"
docker exec goshop-postgres psql -U commerce -d commerce -c "SELECT * FROM gocommerce_migrations;"

# View seeded data
docker exec goshop-postgres psql -U commerce -d commerce -c "SELECT id, name FROM products LIMIT 5;"
docker exec goshop-postgres psql -U commerce -d commerce -c "SELECT id, name FROM categories LIMIT 5;"
docker exec goshop-postgres psql -U commerce -d commerce -c "SELECT id, name FROM brands LIMIT 5;"
```

### Seeding Data

Sample data is provided through the `gocommerce/migrations` package. Set `SEED_DB=true` in your `.env` file to populate the database with:

- **Products**: Various electronics (laptops, phones, tablets)
- **Categories**: Electronics, Computers, Accessories, Audio, Storage
- **Brands**: Apple, Dell, Lenovo, HP, Samsung
- **Variants**: Product variations (sizes, colors)

## Configuration Reference

All configuration is done via environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | 8080 | No |
| `DB_DRIVER` | Database driver | postgres | Yes |
| `DB_DSN` | Database connection string | - | Yes |
| `JWT_SECRET` | JWT signing key (min 32 chars) | - | Yes |
| `JWT_ACCESS_TOKEN_EXPIRY` | Access token lifetime | 15m | No |
| `JWT_REFRESH_TOKEN_EXPIRY` | Refresh token lifetime | 168h | No |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | - | No |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | - | No |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | http://localhost:8080/api/v1/auth/google/callback | No |
| `SEED_DB` | Seed database with sample data | false | No |

## Google OAuth Setup

The API supports **Sign in with Google** for seamless user authentication. See [GOOGLE_OAUTH.md](GOOGLE_OAUTH.md) for complete setup guide.

### Quick Setup

1. **Get Google OAuth Credentials**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create/select a project
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`

2. **Configure Environment Variables**
   ```env
   GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=your-client-secret
   GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
   ```

3. **Test OAuth Flow**
   - Open `test-oauth.html` in your browser
   - Click "Sign in with Google"
   - Authorize the application
   - Receive JWT tokens

### OAuth Endpoints

- `GET /api/v1/auth/google` - Get authorization URL
- `GET /api/v1/auth/google/callback` - OAuth callback handler

For detailed integration examples and troubleshooting, see [GOOGLE_OAUTH.md](GOOGLE_OAUTH.md).

## Extending the Application

### Adding New Endpoints
1. Create handler in `internal/http/handlers/`
2. Add route in `internal/http/server.go`
3. Apply authentication middleware if needed

### Adding Inventory Management
Database tables for inventory are already created by gocommerce v0.0.5 migrations (suppliers, inventory_levels, inventory_activities). To enable inventory features:
1. Create repository implementations for inventory tables
2. Implement the `inventory.Service` interface from gocommerce
3. Inject the service into cart and order services (currently `nil`)

### Adding Payment Processing
Implement the `payments.Gateway` interface from gocommerce and inject it into the order service.

### Adding Shipping Calculators
Implement the `shipping.RateCalculator` interface from gocommerce and inject it into the pricing service.

## License

MIT

## Acknowledgments

- **gocommerce** - E-commerce domain library
- **goauthx** - Authentication & authorization library
- **Gin** - HTTP web framework
- **GORM** - ORM library
