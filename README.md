# Code structure
* `client/client.go` generates and ingests metrics data for the past 5 minutes with 60 milliseconds interval
* `migrations/1_metrics.up.sql` contains schema table for metrics
* `pkg/api.api.go` contains API to return metrics data for any time range
* `pkg/db/db.go` starts PostgreSQL db, creates connection and if necessary deploys schema
* `pkg/db/models/metrics.go` provides API to insert metrics data to DB and get metrics data from DB
* `server/main.go` starts DB, creates a route and listen to requests

# Build
### Generate and ingest data 
`go build -o ./bin ./client`

### API to query data
`go build -o ./bin ./server`

# Run locally (Linux)
* Install docker and docker compose

* Set up Postgres DB

`docker compose up db -d` 

* Set env variables

`export DATABASE_URL="postgresql://postgres:postgres@localhost/postgres?sslmode=disable"`

`export PORT=8080`

* Run client to generate and ingest data

`go run ./client/client.go`

* Run server

`go run ./server/main.go`

* Request data up to the current date use the next curl command

`curl "http://localhost:8080/metrics?from_ts=0&to_ts=$(date +%s)"`


# Test
### To run end-to-end test
`docker compose up --build --exit-code-from test --abort-on-container-exit`

### To clean database
`docker compose down -v --remove-orphans`