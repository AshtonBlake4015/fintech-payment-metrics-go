package main

import (
	"testing"
	"time"
)

func TestPaymentMetricShape(t *testing.T) {
	metric := map[string]any{"type": "counter", "name": "payments.authorized", "value": 1, "tags": map[string]string{"method": "card"}}
	if metric["type"] != "counter" || metric["name"] != "payments.authorized" || metric["value"] != 1 {
		t.Fatal("unexpected counter shape")
	}
	if 184*time.Millisecond > time.Second {
		t.Fatal("sample latency outside test range")
	}
}
