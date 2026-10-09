#!/usr/bin/env bash
set -euo pipefail

cd -- "$(dirname -- "${BASH_SOURCE[0]}")"

minikube kubectl -- apply -f k8s/namespace.yaml
minikube kubectl -- apply -f k8s/configmap.yaml -f k8s/secret.yaml -f k8s/service.yaml -f k8s/postgres.yaml
minikube kubectl -- rollout status statefulset/traci-postgres -n traci --timeout=120s

minikube kubectl -- exec -n traci traci-postgres-0 -- \
  sh -c 'for attempt in $(seq 1 60); do pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" && exit 0; sleep 2; done; exit 1'

minikube kubectl -- delete job traci -n traci --ignore-not-found
minikube kubectl -- apply -f k8s/migrate.yaml
minikube kubectl -- wait -n traci --for=condition=complete job/traci --timeout=120s

minikube kubectl -- apply -f k8s/backend.yaml -f k8s/frontend.yaml -f k8s/ingress.yaml
