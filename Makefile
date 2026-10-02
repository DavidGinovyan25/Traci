.DEFAULT_GOAL := build

.PHONY: generate build target test test-unit test-race run docker-up docker-down

generate:
	cd backend && go generate ./internal/api
	cd backend && sqlc generate
	cd backend && go run ./internal/api/stripcomments internal/gen/db/db.go internal/gen/db/models.go internal/gen/db/querier.go internal/gen/db/user.sql.go internal/gen/db/movie.sql.go internal/gen/db/collection.sql.go

build:
	cd backend && go build -o bin/traci ./cmd/traci

target: build

test:
	cd backend && go test ./...
	cd backend && go vet ./...

test-unit:
	cd backend && go test -short ./...

test-race:
	cd backend && go test -race ./...

run:
	cd backend && go run ./cmd/traci

docker-up:
	docker compose --env-file backend/.env up --build -d --wait

docker-down:
	docker compose --env-file backend/.env down
