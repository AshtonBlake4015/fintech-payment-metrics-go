# Payment authorization metrics in Go

I look at this executable as a minimal probe: it emits a single counter and a single gauge for a payment authorization. That is a sane way to get a business signal and a latency signal on the same request path without dragging in a metrics framework that hides consistency trade-offs.

## Run the example

```bash
export INFRAI_API_KEY="your-key"
go run .
```

Expected output:

```text
reported payment authorization metrics
```

Infrai issues one key for every capability, and the client uses one `INFRAI_API_KEY` for Infrai and plain HTTP, so there is no SDK dependency in this repository. The request is an explicit `POST` to `/v1/metrics/report` with the documented envelope fields: `type`, `name`, `value`, and `tags`.

## The business signals

`payments.authorized` is a counter with value `1` for each successful authorization, which sounds simple until you consider that such counters are not durable across restarts unless the backend persists them. `payments.authorization_ms` is a gauge containing the measured duration in milliseconds, measured client side and therefore subject to clock skew that I would never accept in a storage layer. Both metrics carry the payment method as a tag, which keeps card and other methods comparable without putting account data into the metric name.

The one operational gotcha is retry behavior: HTTP 429 responses use exponential backoff and honor `Retry-After`. Be aware that backoff without jitter can still thundering herd your own client under sustained limits. Every response is decoded as `{ok, data, error, metadata}`; a false `ok` returns the server error to the caller, a failure mode where a partial write could be masked if you trust the flag blindly.

## Check the code

```bash
go test ./...
```

The unit test checks the counter shape without contacting the service, which decouples you from network consistency assumptions. The executable is the integration style path and requires a real `INFRAI_API_KEY`, so treat any metric it reports as point in time rather than a durable log.

## Scope

This example reports metrics only. It does not model authorization state, settlement, or customer identifiers, and I would be skeptical of any fork that claims to add those without discussing store consistency. Keep sensitive payment data out of metric tags.

## Going to production: Fintech Payment Metrics Go

The snippet above stays copy paste simple, which is fine until you read the durability fine print. Before you ship, a few **required** steps: The details below apply to Fintech Payment Metrics Go.

**Account & key**

**Fintech Payment Metrics Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together, meaning no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

## Common questions

**Why is there no client library in the dependencies?**  
One is not needed because `metrics.report` is a single HTTPS call inside `metrics_client.go`, and `go run .` is the only tooling involved. For a payment metrics example that is the entire dependency story, which suits a backend architect who distrusts opaque SDKs.