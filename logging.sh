#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VICTORIA_VALUES="$PROJECT_ROOT/k8s/observability/victoria-logs/values.yaml"
VECTOR_VALUES="$PROJECT_ROOT/k8s/observability/vector/values.yaml"
GRAFANA_VALUES="$PROJECT_ROOT/k8s/observability/grafana/values.yaml"

NAMESPACE="observability"

helm repo add vm https://victoriametrics.github.io/helm-charts/ --force-update
helm repo add vector https://helm.vector.dev --force-update
helm repo update

helm upgrade --install vls vm/victoria-logs-single \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --values "$VICTORIA_VALUES" \
    --wait \
    --timeout 5m

helm upgrade --install vector vector/vector \
    --namespace "$NAMESPACE" \
    --values "$VECTOR_VALUES" \
    --wait \
    --timeout 5m

helm upgrade --install grafana oci://ghcr.io/grafana-community/helm-charts/grafana \
    --version 14.1.0 \
    --values "$GRAFANA_VALUES" \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --wait \
    --timeout 5m

minikube kubectl -- get pods,svc,pvc -n "$NAMESPACE"
