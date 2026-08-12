package main

import (
	"fmt"
	"time"
)

func reportPaymentMetrics(client *MetricsClient, method string, amountCents int64, elapsed time.Duration) error {
	tags := map[string]string{"method": method}
	if err := client.Report(map[string]any{"type": "counter", "name": "payments.authorized", "value": 1, "tags": tags}); err != nil {
		return err
	}
	return client.Report(map[string]any{"type": "gauge", "name": "payments.authorization_ms", "value": elapsed.Milliseconds(), "tags": tags})
}

func main() {
	client, err := NewMetricsClient()
	if err != nil {
		panic(err)
	}
	if err := reportPaymentMetrics(client, "card", 12500, 184*time.Millisecond); err != nil {
		panic(err)
	}
	fmt.Println("reported payment authorization metrics")
}
