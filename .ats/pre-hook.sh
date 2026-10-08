#!/usr/bin/env bash
# Deploys the chart under test the way it runs in production: next to a Gel
# (EdgeDB) server, which the operator connects to at startup and exits without.
# app-test-suite runs this before the tests, with its own chart deploy skipped
# (.ats/main.yaml), so the server and the secrets the chart mounts exist first.
set -euo pipefail

: "${ATS_CHART_PATH:?}"
namespace=default
here=$(dirname "$0")

kubectl -n "$namespace" create secret generic edgedb-server-password \
  --from-literal=password="$(head -c 24 /dev/urandom | base64 | tr -d '/+=')"
kubectl -n "$namespace" apply -f "$here/edgedb.yaml"
kubectl -n "$namespace" rollout status deployment/edgedb --timeout=10m

# The server generates a self-signed certificate on first start; the operator
# trusts it as its CA.
kubectl -n "$namespace" exec deployment/edgedb -- cat /var/lib/gel/data/edbtlscert.pem > "$here/edgedb-ca.crt"
kubectl -n "$namespace" create secret generic edgedb-tls --from-file=ca.crt="$here/edgedb-ca.crt"
rm -f "$here/edgedb-ca.crt"

helm upgrade --install policy-meta-operator "$ATS_CHART_PATH" \
  --namespace "$namespace" --wait --timeout 10m
