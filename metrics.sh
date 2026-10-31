#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VICTORIA_METRICS_VALUES="$PROJECT_ROOT/k8s/observability/victoria-metrics.yaml"

NAMESPACE="observability"

helm repo add vm https://victoriametrics.github.io/helm-charts/ --force-update
helm repo update

helm upgrade --install vmsingle vm/victoria-metrics-single \
    --values "$VICTORIA_METRICS_VALUES" \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --wait \
    --timeout 5m