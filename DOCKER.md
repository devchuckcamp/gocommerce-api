# Docker Deployment Guide

This guide covers how to run the GoShop E-Commerce API using Docker and Docker Compose.

## Quick Start

### 1. Ensure Environment Variables are Set

Make sure your `.env` file has the required variables:

```env
# Database
DB_USER=commerce
DB_PASSWORD=commerce
DB_NAME=commerce
DB_PORT=5432
DB_DRIVER=postgres

# JWT (REQUIRED - must be at least 32 characters)
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-change-this

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

# Seeding
SEED_DB=true
```

### 2. Build and Start All Services

```bash
docker-compose up --build
```

This will:
- Build the Go API Docker image
- Start PostgreSQL database
- Start the API service
- Run migrations automatically
- Seed database if `SEED_DB=true`

### 3. Access the API

The API will be available at:
- **API Base URL**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/health

## Docker Compose Services

### PostgreSQL Database
- **Image**: postgres:16-alpine
- **Container**: goshop-postgres
- **Port**: 5432 (configurable via `DB_PORT`)
- **Volume**: postgres_data (persistent storage)
- **Health Check**: Automated readiness check

### API Service
- **Image**: Built from Dockerfile
- **Container**: goshop-api
- **Port**: 8080 (configurable via `PORT`)
- **Health Check**: HTTP check on /health endpoint
- **Depends On**: PostgreSQL (waits for database to be healthy)

## Docker Commands

### Start Services
```bash
# Start in foreground (see logs)
docker-compose up

# Start in background (detached)
docker-compose up -d

# Rebuild and start
docker-compose up --build
```

### Stop Services
```bash
# Stop services (preserves data)
docker-compose down

# Stop and remove volumes (deletes database data)
docker-compose down -v
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f postgres
```

### Rebuild API
```bash
# Rebuild just the API service
docker-compose build api

# Force rebuild without cache
docker-compose build --no-cache api
```

### Execute Commands in Containers
```bash
# Access PostgreSQL
docker-compose exec postgres psql -U commerce -d commerce

# Access API container shell
docker-compose exec api sh

# Run database queries
docker-compose exec postgres psql -U commerce -d commerce -c "SELECT * FROM users LIMIT 5;"
```

## Environment Variables

The API container receives environment variables from your `.env` file. Key variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | 8080 |
| `DB_HOST` | Database hostname | postgres |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database username | commerce |
| `DB_PASSWORD` | Database password | commerce |
| `DB_NAME` | Database name | commerce |
| `JWT_SECRET` | JWT signing key (REQUIRED) | - |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth secret | - |
| `SEED_DB` | Seed with sample data | false |

## Production Deployment

### 1. Update Environment Variables

For production, update your `.env` or set environment variables:

```env
# Use strong, random JWT secret
JWT_SECRET=$(openssl rand -base64 48)

# Use strong database password
DB_PASSWORD=$(openssl rand -base64 32)

# Update OAuth redirect URL for your domain
GOOGLE_REDIRECT_URL=https://your-domain.com/api/v1/auth/google/callback

# Disable seeding in production
SEED_DB=false
```

### 2. Build Production Image

```bash
# Build with production tag
docker-compose build

# Or build manually
docker build -t goshop-api:1.0.0 .
```

### 3. Push to Container Registry

```bash
# Tag for your registry
docker tag goshop-api:1.0.0 your-registry.com/goshop-api:1.0.0

# Push to registry
docker push your-registry.com/goshop-api:1.0.0
```

### 4. Deploy with Docker Compose

Update `docker-compose.yml` to use the registry image:

```yaml
services:
  api:
    image: your-registry.com/goshop-api:1.0.0
    # ... rest of config
```

## Dockerfile Details

The application uses a **multi-stage build** for optimal image size:

### Build Stage
- Base: `golang:1.23-alpine`
- Installs build dependencies
- Downloads Go modules
- Compiles static binary with optimizations
- Result: ~15MB binary

### Runtime Stage
- Base: `alpine:3.19`
- Minimal dependencies (ca-certificates, tzdata)
- Non-root user for security
- Only includes compiled binary
- Final image: ~25MB

### Security Features
- ✅ Non-root user (appuser:1000)
- ✅ Minimal attack surface (Alpine Linux)
- ✅ No build tools in runtime image
- ✅ Read-only filesystem compatible
- ✅ Health checks enabled

## Networking

Services communicate through `goshop-network` bridge network:
- API connects to postgres via hostname `postgres`
- External access via exposed ports
- Database not exposed to host by default (more secure)

## Health Checks

### PostgreSQL
```bash
pg_isready -U commerce -d commerce
```
- Interval: 10s
- Timeout: 5s
- Retries: 5

### API
```bash
wget --spider http://localhost:8080/health
```
- Interval: 30s
- Timeout: 3s
- Start Period: 10s (allows time for migrations)
- Retries: 3

## Troubleshooting

### API Won't Start

Check logs:
```bash
docker-compose logs api
```

Common issues:
- JWT_SECRET not set or too short
- Database not ready (wait for health check)
- Port 8080 already in use

### Database Connection Failed

Check database health:
```bash
docker-compose ps
docker-compose logs postgres
```

Test connection:
```bash
docker-compose exec postgres psql -U commerce -d commerce -c "SELECT 1;"
```

### Migrations Failed

View migration status:
```bash
docker-compose exec postgres psql -U commerce -d commerce -c "SELECT * FROM schema_migrations;"
```

Reset database (WARNING: deletes all data):
```bash
docker-compose down -v
docker-compose up --build
```

### OAuth Not Working

Ensure redirect URL is correct:
- For Docker: Use `http://localhost:8080/api/v1/auth/google/callback`
- For production: Update to your domain
- Must match exactly in Google Cloud Console

### Port Already in Use

Change the port in `.env`:
```env
PORT=8081
```

Then rebuild:
```bash
docker-compose down
docker-compose up --build
```

## Monitoring

### Container Stats
```bash
docker stats goshop-api goshop-postgres
```

### Resource Usage
```bash
docker-compose ps
docker system df
```

### Database Size
```bash
docker-compose exec postgres psql -U commerce -d commerce -c "
SELECT pg_size_pretty(pg_database_size('commerce')) AS db_size;
"
```

## Backup and Restore

### Backup Database
```bash
docker-compose exec postgres pg_dump -U commerce commerce > backup.sql
```

### Restore Database
```bash
docker-compose exec -T postgres psql -U commerce commerce < backup.sql
```

### Backup Volume
```bash
docker run --rm \
  -v goshop_postgres_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres-backup.tar.gz /data
```

## Scaling

### Horizontal Scaling (Multiple API Instances)

Update `docker-compose.yml`:
```yaml
services:
  api:
    # ... config
    deploy:
      replicas: 3
```

Add a load balancer (nginx, traefik, etc.) in front of API instances.

### Vertical Scaling (Resource Limits)

```yaml
services:
  api:
    # ... config
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build and Push Docker Image

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: docker build -t goshop-api:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo ${{ secrets.REGISTRY_TOKEN }} | docker login -u ${{ secrets.REGISTRY_USER }} --password-stdin
          docker push goshop-api:${{ github.sha }}
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Alpine Linux](https://alpinelinux.org/)
