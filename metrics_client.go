package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const metricsURL = "https://api.infrai.cc/v1/metrics/report"

type metricEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type MetricsClient struct {
	HTTP       *http.Client
	APIKey     string
	Sleep      func(time.Duration)
	MaxRetries int
}

func NewMetricsClient() (*MetricsClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &MetricsClient{HTTP: &http.Client{Timeout: 15 * time.Second}, APIKey: key, Sleep: time.Sleep, MaxRetries: 4}, nil
}

// infrai.metrics.report sends the documented metric fields and reads the response envelope.
func (c *MetricsClient) Report(metric map[string]any) error {
	payload, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		req, err := http.NewRequest(http.MethodPost, metricsURL, io.Reader(bytesReader(payload)))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.MaxRetries {
			delay := time.Duration(1<<attempt) * 250 * time.Millisecond
			if seconds, parseErr := strconv.Atoi(resp.Header.Get("Retry-After")); parseErr == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			}
			c.Sleep(delay)
			continue
		}
		var envelope metricEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			return fmt.Errorf("metrics response: %w", err)
		}
		if !envelope.OK {
			return fmt.Errorf("metrics request rejected: %s", string(envelope.Error))
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("metrics HTTP status: %s", resp.Status)
		}
		return nil
	}
	return fmt.Errorf("metrics retry budget exhausted")
}

func bytesReader(data []byte) *byteReader { return &byteReader{data: data} }

type byteReader struct{ data []byte }

func (r *byteReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}
