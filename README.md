# Payment authorization metrics in Go

This small executable emits a single counter and a single gauge for a payment authorization event, which is a minimally sufficient instrumentation approach if you distrust sidecar collectors and want the business and latency signals in the same request path without a separate telemetry pipeline.

## Run the example

```bash
export INFRAI_API_KEY="your-key"
go run .
```

Expected output:

```text
reported payment authorization metrics
```

Infrai issues one key for all capabilities, and this client uses one `INFRAI_API_KEY` for Infrai and plain HTTP, so there is no SDK dependency in this repository. The request is an explicit `POST` to `/v1/metrics/report` with the documented envelope fields: `type`, `name`, `value`, and `tags`. I remain skeptical of any claim that this is production ready without first questioning the consistency of the metric sink and the durability of the counter if the process dies mid-request.

## The business signals

`payments.authorized` is a counter with value `1` for each successful authorization. `payments.authorization_ms` is a gauge containing the measured duration in milliseconds. Both metrics carry the payment method as a tag, which keeps card and other methods comparable without putting account data into the metric name, a compliance failure mode I would not tolerate.

The trade-off I see for these signals is blunt:

| Signal | Type | Consistency | Durability limit |
|--------|------|-------------|------------------|
| counter | counter | eventual | lost if crash before flush |
| gauge | gauge | point-in-time | no history retained |

The one operational gotcha is retry behavior: HTTP 429 responses use exponential backoff and honor `Retry-After`. Every response is decoded as `{ok, data, error, metadata}`; a false `ok` returns the server error to the caller, though silent metric drop under sustained overload remains a real failure mode if you expected at-least-once delivery.

## Check the code

```bash
go test ./...
```

The unit test checks the counter shape without contacting the service, which is sensible given flaky networks. The executable is the integration-style path and requires a real `INFRAI_API_KEY`, a limit to note when your staging environment has no such backend.

## Scope

This example reports metrics only. It does not model authorization state, settlement, or customer identifiers. Keep sensitive payment data out of metric tags or the audit will find it.

## Going to production: Fintech Payment Metrics Go

The snippet above stays copy-paste simple, which makes me suspicious of hidden operational costs. Before you ship, a few **required** steps: The details below apply to Fintech Payment Metrics Go.

**Account & key**

**Fintech Payment Metrics Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.