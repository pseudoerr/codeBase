# Auth Service

JWT-based authentication service for interactive learning Platform (CodeBase) with secure user registration, login, and token management.

## Features

- **User Registration & Login** with email/password
- **JWT Authentication** with access/refresh token pattern
- **Secure Password Hashing** using bcrypt
- **Token Management** with database-stored refresh tokens
- **Input Validation** with comprehensive error handling
- **Structured Logging** with slog
- **Graceful Shutdown** support
- **Health Check** endpoint
- **Docker Support** with docker-compose
- **Database Migrations** with golang-migrate
- **Unit Tests** with mocks

## API Endpoints

### Public Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh access token
- `GET /health` - Health check

### Protected Endpoints (require JWT)
- `GET /auth/me` - Get user profile
- `POST /auth/logout` - Logout user

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- golang-migrate CLI tool

### Local Development

1. **Clone and Setup**
   ```bash
   git clone <repository>
   cd auth-service
   make dev-setup
   ```

2. **Run Migrations**
   ```bash
   make migrate-up
   ```

3. **Start Service**
   ```bash
   make run
   ```

### Docker Development

1. **Start All Services**
   ```bash
   make docker-up
   ```

2. **Check Logs**
   ```bash
   make docker-logs
   ```

## Usage Examples

### Register User
```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "SecurePass123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Get Profile
```bash
curl -X GET http://localhost:8081/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Refresh Token
```bash
curl -X POST http://localhost:8081/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Configuration

Environment variables (see `.env.example`):

- `PORT` - Server port (default: 8081)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret (change in production!)
- `JWT_ACCESS_TTL` - Access token TTL (default: 15m)
- `JWT_REFRESH_TTL` - Refresh token TTL (default: 168h)
- `BCRYPT_COST` - Bcrypt hashing cost (default: 12)

## Security Features

- **Password Requirements**: Minimum 8 characters with uppercase, lowercase, and numbers
- **JWT Security**: Short-lived access tokens (15 min) with secure refresh mechanism
- **Password Hashing**: bcrypt with configurable cost
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Configurable for production
- **Rate Limiting**: Ready for implementation

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Security scan
make security
```

## Database Schema

### Users Table
- `id` - Primary key
- `email` - Unique email address
- `username` - Unique username
- `password_hash` - Bcrypt hashed password
- `created_at`, `updated_at` - Timestamps

### Refresh Tokens Table
- `id` - Primary key
- `user_id` - Foreign key to users
- `token` - Unique refresh token
- `expires_at` - Token expiration
- `created_at` - Creation timestamp

## Integration with Other Services

This service provides JWT middleware that can be imported by other services:

```go
import "auth-service/internal/middleware"

// Use JWT middleware
router.Use(middleware.JWTMiddleware(jwtSecret))

// Access user info from request headers
userID := r.Header.Get("X-User-ID")
email := r.Header.Get("X-User-Email")
username := r.Header.Get("X-User-Username")
```

## Production Checklist

- [ ] Change JWT_SECRET to secure random value
- [ ] Configure CORS for your frontend domain
- [ ] Set up proper SSL/TLS certificates
- [ ] Configure rate limiting
- [ ] Set up monitoring and alerting
- [ ] Configure log aggregation
- [ ] Set up database backups
- [ ] Configure environment-specific configs
- [ ] Set up CI/CD pipeline
- [ ] Security audit and penetration testing

## TODO for Production

- [ ] Rate limiting middleware implementation
- [ ] Password reset functionality
- [ ] Email verification
- [ ] OAuth integration (Google, GitHub)
- [ ] Account lockout after failed attempts
- [ ] Audit logging for security events
- [ ] Metrics collection (Prometheus)
- [ ] Health check with database connectivity
- [ ] Configuration management (Vault)
- [ ] Load testing and performance optimization