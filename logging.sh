#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VICTORIA_VALUES="$PROJECT_ROOT/k8s/observability/victoria-logs.yaml"
VECTOR_VALUES="$PROJECT_ROOT/k8s/observability/vector-values.yaml"
GRAFANA_VALUES="$PROJECT_ROOT/k8s/observability/grafana.yaml"
VMALERT_VALUES="$PROJECT_ROOT/k8s/observability/vmalert.yaml"
ALERTMANAGER_VALUES="$PROJECT_ROOT/k8s/observability/alertmanager.yaml"

NAMESPACE="observability"

helm repo add vm https://victoriametrics.github.io/helm-charts/ --force-update
helm repo add vector https://helm.vector.dev --force-update
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
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

helm upgrade --install alertmanager prometheus-community/alertmanager \
    --values "$ALERTMANAGER_VALUES" \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --wait \
    --timeout 5m

helm upgrade --install vma vm/victoria-metrics-alert \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --values "$VMALERT_VALUES" \
    --wait \
    --timeout 5m

minikube kubectl -- get pods,svc,pvc -n "$NAMESPACE"
