#  Mission Service API

REST API for the gamified coding platform.

---

## 🚀 Key Features

- CRUD API for `/missions`
- Gamification: `/profile` that supports **missions** and **badges** 
- PostgreSQL + simple migrations (`golang-migrate`)
- Middleware: structured logging with slog, CORS, panic/recovery
- Unit tests (`httptest`)
- Clean Architecture: `handlers`, `repository`, `service`, `config`

---

##  Project Structure

```

.
├── cmd/main.go              # Entrypoint
├── config/                  # Env-variables loader 
├── internal/http/           # Handlers, routers, middleware
├── migrations/              # Sql-files for migrations
├── models/                  # DTO-models (Mission, Profile)
├── repository/              # PostgreSQL-repo
├── service/                 # Business logic

````

---

## Installation 

1. Install dependecies:

```bash
go mod tidy
````

2. Prepare PostgreSQL:

```bash
docker run --name=missions-db \
  -e POSTGRES_PASSWORD=missions \
  -p 5436:5432 \
  -d postgres
```

3. Create`.env` file:

```
PORT=8080
DATABASE_URL=postgres://postgres:missions@localhost:5436/postgres?sslmode=disable
```

4. Migrate:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

5. Launch application:

```bash
go run ./cmd
```


##  Examples of simple CURL-requests

 Create new mission:

```bash
curl -X POST http://localhost:8080/missions \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello, World!", "points": 100}'
```
 
Get all missions:

```bash
curl http://localhost:8080/missions
```

Get profile:

```bash
curl http://localhost:8080/profile
```

Delete missions:

```bash
curl -X DELETE http://localhost:8080/missions/1
```

---

###  ToDo 

* [ ] Dockerfile + docker-compose
* [ ] CI (Go test + lint)
* [ ] Swagger/OpenAPI
* [ ] Frontend-interface