.DEFAULT_GOAL := build

.PHONY: seed clear

.PHONY: generate build target test test-unit test-race run docker-up docker-down test-frontend test-e2e create-admin images-build minikube-load minikube-deploy-app minikube-deploy-observability minikube-delete-observability

	
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

seed:
	minikube kubectl -- exec -i -n traci traci-postgres-0 \
		-c traci-postgres -- \
		psql -X -v ON_ERROR_STOP=1 -U traci -d traci \
		< backend/db/postgres/seed.sql

clear:
	minikube kubectl -- exec -i -n traci traci-postgres-0 \
		-c traci-postgres -- \
		psql -X -v ON_ERROR_STOP=1 -U traci -d traci \
		< backend/db/postgres/clear.sql

images-build:
	docker build -f backend/cmd/traci/Dockerfile -t traci-backend:v1 backend
	docker build -t traci-frontend:v1 frontend
	docker build -t traci-prober:v1 prober

minikube-load: images-build
	minikube image load traci-backend:v1
	minikube image load traci-frontend:v1
	minikube image load traci-prober:v1

minikube-deploy-app:
	./deploy.sh

minikube-deploy-observability:
	./deploy_observability.sh

minikube-delete-observability:
	helm uninstall alertmanager -n observability
	helm uninstall grafana -n observability
	helm uninstall vector -n observability
	helm uninstall vls -n observability
	helm uninstall vma -n observability
	helm uninstall vmsingle -n observability
