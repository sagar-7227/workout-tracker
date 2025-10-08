# Workout Tracker API (Go)


## Quick start
```bash
cp .env.example .env
# edit DB_DSN if needed


# start postgres (optional)
docker compose up -d


# run server
go run ./cmd/server
```


### Auth
- `POST /api/auth/signup` {name,email,password}
- `POST /api/auth/login` {email,password} → `{ token }`


### Exercises
- `GET /api/exercises` (auth)


### Workouts
- `POST /api/workouts` create workout + exercises
- `GET /api/workouts` list (sorted by `scheduled_at` asc)
- `GET /api/workouts/{id}` get
- `PUT /api/workouts/{id}` update (replaces exercises if provided)
- `DELETE /api/workouts/{id}` delete


### Reports
- `GET /api/reports` summary


### OpenAPI (Swagger)
Install generator and generate docs from annotations:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency --parseInternal --dir ./cmd/server,./controllers --output ./docs
```
Serve at `/swagger/index.html` (in non-prod).


### Testing
```bash
go test ./...
```