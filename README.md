# Test Task

A microservices system comprising two Go services and infrastructure components (PostgreSQL, RabbitMQ).

## Project Structure

- `products/`: REST API for product management (Gin, pgx, RabbitMQ Producer, Prometheus Metrics).
- `notification-service/`: Background worker for logging product events (RabbitMQ Consumer).
- `docker-compose.yml`: Orchestrates infrastructure and microservices.
- `scripts/`: Useful SQL and shell scripts for seeding and testing.

## Requirements

- Docker and Docker Compose
- Go 1.25

---

## How to Run

### Using Docker Compose (Recommended)

Run everything with a single command from the project root:

```bash
docker compose up --build
```

This will:

1. Start **PostgreSQL** (Port: 5432)
2. Start **RabbitMQ** (Port: 5672, Management: 15672)
3. Run **Products Service** (Port: 8080)
4. Start **Notifications Service** (Background consumer)

Database migrations will be applied automatically on Products Service startup.

---

## How to Test

### Run Unit Tests

To run unit tests for the **Products Service**, ensuring business logic for creation and deletion works correctly:

```bash
cd products
go test ./...
```

_(Or use `go test -v ./internal/service` to see detailed execution for the service layer with mocks.)_

---

## API Endpoints & Usage

### Products Service

- `GET /ping`: Health check.
- `GET /metrics`: Prometheus metrics (counters for create/delete).
- `POST /products`: Create a product.
- `DELETE /products/:id`: Delete a product.
- `GET /products/list?limit=10&page=1`: Get products with pagination (page-based).

### Examples (using curl)

1. **Create a Product**:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Coffee Machine", "price": 199.50}'
```

2. **List Products (Page 1)**:

```bash
curl "http://localhost:8080/products/list?limit=10&page=1"
```

3. **Check Prometheus Metrics**:

```bash
curl http://localhost:8080/metrics | grep products_
```

---

## Database Seeding

If you want to quickly populate the database with test data:

1. Connect to your PostgreSQL instance (localhost:5432).
2. Use the SQL in `scripts/seed.sql`:

```bash
# Example if you have psql installed:
psql -U user -d catalog_db -h localhost -f scripts/seed.sql
```
