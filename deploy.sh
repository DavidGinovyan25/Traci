#!/usr/bin/env bash
set -euo pipefail

cd -- "$(dirname -- "${BASH_SOURCE[0]}")"

minikube kubectl -- apply -f k8s/app/namespace.yaml
minikube kubectl -- apply -f k8s/app/configmap.yaml -f k8s/app/service.yaml
./setup_postgres_exporter.sh

minikube kubectl -- exec -n traci traci-postgres-0 -- \
  sh -c 'for attempt in $(seq 1 60); do pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" && exit 0; sleep 2; done; exit 1'

minikube kubectl -- delete job traci -n traci --ignore-not-found
minikube kubectl -- apply -f k8s/app/migrate.yaml
minikube kubectl -- wait -n traci --for=condition=complete job/traci --timeout=120s

minikube kubectl -- apply -f k8s/app/backend.yaml -f k8s/app/frontend.yaml -f k8s/app/ingress.yaml
