#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$PROJECT_ROOT/backend/.env"

minikube kubectl -- create secret generic traci-secret \
  -n traci \
  --from-env-file=<(grep -E '^(POSTGRES_PASSWORD|JWT_SECRET)=' "$ENV_FILE") \
  --dry-run=client -o yaml |
  minikube kubectl -- apply -f -

minikube kubectl -- create secret generic traci-exporter-secret \
  -n traci \
  --from-env-file=<(grep '^DATA_SOURCE_PASS=' "$ENV_FILE") \
  --dry-run=client -o yaml |
  minikube kubectl -- apply -f -

minikube kubectl -- apply -f "$PROJECT_ROOT/k8s/app/postgres.yaml"
minikube kubectl -- rollout status statefulset/traci-postgres -n traci --timeout=180s

minikube kubectl -- exec -i -n traci traci-postgres-0 \
  -c traci-postgres -- \
  psql -X -v ON_ERROR_STOP=1 -U traci -d traci \
  < "$PROJECT_ROOT/backend/db/postgres/new_monitor_user.sql"
