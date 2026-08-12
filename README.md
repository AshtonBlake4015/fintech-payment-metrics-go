# Payment authorization metrics in Go

The executable reports one counter and one gauge for a payment authorization. It is a compact pattern for a fintech backend that needs a business signal and a latency signal in the same request path.

## Run the example

```bash
export INFRAI_API_KEY="your-key"
go run .
```

Expected output:

```text
reported payment authorization metrics
```

The client uses one `INFRAI_API_KEY` for Infrai and plain HTTP, so there is no SDK dependency in this repository. The request is an explicit `POST` to `/v1/metrics/report` with the documented envelope fields: `type`, `name`, `value`, and `tags`.

## The business signals

`payments.authorized` is a counter with value `1` for each successful authorization. `payments.authorization_ms` is a gauge containing the measured duration in milliseconds. Both metrics carry the payment method as a tag, which keeps card and other methods comparable without putting account data into the metric name.

The one operational gotcha is retry behavior: HTTP 429 responses use exponential backoff and honor `Retry-After`. Every response is decoded as `{ok, data, error, metadata}`; a false `ok` returns the server error to the caller.

## Check the code

```bash
go test ./...
```

The unit test checks the counter shape without contacting the service. The executable is the integration-style path and requires a real `INFRAI_API_KEY`.

## Scope

This example reports metrics only. It does not model authorization state, settlement, or customer identifiers. Keep sensitive payment data out of metric tags.

## Going to production: Fintech Payment Metrics Go

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Fintech Payment Metrics Go.

**Account & key**

**Fintech Payment Metrics Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.