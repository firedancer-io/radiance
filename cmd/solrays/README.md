# solrays

An RPC proxy that exports call latency of a Solana RPC node.

## Metrics

- `solrays_forwarded_duration_seconds` - an histogram of the duration of the calls forwarded to the Solana RPC node, per call.
- `solrays_requests_total` - Count of total requests per call.
- `solrays_request_errors_total` - Count of total requests that failed, e.g. malformed request or proxy can't reach the Solana RPC node.
- `solrays_requests_status_total` - Count of response HTTP status codes, per status code.

## Run

```
:; _bin/solrays -h
```

Key flags are:

- `backend` - the Solana RPC endpoint to proxy and monitor. Defaults to `http://localhost:8899`).
- `listen` - address where the RPC metrics and Go profiler endpoints are exposed.
- `tlsHostname` - the optional hostname for this service. Implicitly enables HTTPS.
- `tlsProd` - if HTTPS is enabled, select whether the x509 artifacts are provisioned by the production or the testing Let's Encrypt
  public ACME service. Defaults to `false`, meaning the testing service will be used.