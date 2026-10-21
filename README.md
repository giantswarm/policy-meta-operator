[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/policy-meta-operator/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/policy-meta-operator/tree/main)

# policy-meta-operator

`policy-meta-operator` is a Kubernetes operator that indexes the policy landscape of a cluster into a [Gel](https://github.com/geldata/gel) (formerly EdgeDB) graph database.
It watches enforcing Kyverno `ClusterPolicy` resources and the [Policy API](https://github.com/giantswarm/policy-api) `Policy`, `PolicyConfig`, `PolicyException` and `AutomatedException` resources, mirrors them into the database, and every 30 seconds writes one `PolicyManifest` per policy that aggregates the policy's enforcement mode and every `PolicyException` and `AutomatedException` that targets it.
Those manifests give other Giant Swarm security components a single, queryable view of which workloads are exempt from which policy, without each of them re-reading Kyverno and Policy API objects.

## Installing

The operator is part of Giant Swarm's internal security platform and is not offered as a managed app to customers.
It requires a reachable Gel server (Giant Swarm runs [gel-app](https://github.com/giantswarm/gel-app)) and the Policy API CRDs, which are vendored into the chart from [policy-api](https://github.com/giantswarm/policy-api) via [`vendir.yml`](vendir.yml).
Deploy it like any other Giant Swarm app with a [Flux HelmRelease](https://docs.giantswarm.io/tutorials/fleet-management/app-platform/deploy-app-helmrelease/) or by [creating an App resource](https://docs.giantswarm.io/tutorials/continuous-deployment/apps/add-appcr/).

## Configuring

The only values that normally need changing are the Gel connection settings; the operator exits at startup if it cannot connect.
The password and TLS certificate are read from Secrets in the release namespace.

```yaml
# values.yaml
edgedb:
  host: edgedb
  port: 5656
  user: edgedb
  passwordSecret:
    name: edgedb-server-password
  tlsSecret:
    name: edgedb-tls
```

See our [full reference on how to configure apps](https://docs.giantswarm.io/tutorials/fleet-management/app-platform/app-configuration/) for more details.

## Development

Run `make test-unit` for the Ginkgo unit tests.
The chart tests in `tests/ats` deploy a standalone Gel server first (see [`.ats/edgedb.yaml`](.ats/edgedb.yaml)) because the operator cannot start without one.

## Credit

- [Gel](https://github.com/geldata/gel) and the [edgedb-go](https://github.com/edgedb/edgedb-go) client.
- [Kyverno](https://github.com/kyverno/kyverno), whose `ClusterPolicy` resources are indexed.
