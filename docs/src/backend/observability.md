# Observability

## Prometheus metrics

Prometheus metrics are exposed on `/metrics`.

As we use the default registry of the
[Prometheus GO client library](https://github.com/prometheus/client_golang), Go runtime metrics and
process metrics are also exposed.
