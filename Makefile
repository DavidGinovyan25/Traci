.DEFAULT_GOAL := build

.PHONY: generate build target test test-unit test-race run docker-up docker-down test-frontend test-e2e create-admin images-build minikube-load minikube-deploy

generate:
	cd backend && go generate ./internal/api
	cd backend && go tool sqlc generate

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

test-frontend:
	cd frontend && npm test

test-e2e:
	cd frontend && npm run test:e2e

create-admin:
	minikube kubectl -- exec -it -n traci deployment/traci-backend -- /app/traci create-admin

images-build:
	docker build -f backend/cmd/traci/Dockerfile -t traci-backend:v1 backend
	docker build -t traci-frontend:v1 frontend

minikube-load: images-build
	minikube image load traci-backend:v1
	minikube image load traci-frontend:v1

minikube-deploy:
	./deploy.sh
	minikube kubectl -- apply -f k8s/ingress.yaml
