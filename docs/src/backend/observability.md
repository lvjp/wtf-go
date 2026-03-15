# Observability

## Prometheus metrics

Prometheus metrics are exposed on `/metrics`.

As we use the default registry of the
[Prometheus GO client library](https://github.com/prometheus/client_golang), Go runtime metrics and
process metrics are also exposed.

Specific metrics :

`wtf_go_info`:
:   Information about wtf-go build.  
    Labels:

    * `revision`: the revision identifier for the current commit or checkout.
    * `revision_time`: the modification time associated with vcs.revision, in RFC3339 format.
    * `modified`: true or false indicating whether the source tree had local modifications.

`wtf_go_start_date_timestamp`:
:   The date on which the server started expressed as an UTC Unix timestamp.
