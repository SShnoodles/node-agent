package node

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// AppInfoReporter sends heartbeat and error reports to the center server.
type AppInfoReporter struct {
	cfg       Config
	collector *AppInfoCollector
	client    *http.Client
}

// NewAppInfoReporter creates a new AppInfoReporter.
func NewAppInfoReporter(cfg Config, collector *AppInfoCollector) *AppInfoReporter {
	return &AppInfoReporter{
		cfg:       cfg,
		collector: collector,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// StartHeartbeat sends a heartbeat report every second until ctx is cancelled.
// It is meant to be run in a goroutine.
func (r *AppInfoReporter) StartHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.Report()
		}
	}
}

// Report sends a single heartbeat report. Errors are silently ignored so the
// main application is never affected.
func (r *AppInfoReporter) Report() {
	info, err := r.collector.Collect()
	if err != nil {
		return
	}

	data, err := json.Marshal(info)
	if err != nil {
		return
	}

	r.post(r.cfg.ReportURL, data)
}

// ReportError sends an error report. Returns true if the report was sent
func (r *AppInfoReporter) ReportError(err error) bool {
	if err != nil && err.Error() == "" {
		return true
	}

	info, colErr := r.collector.Collect()
	if colErr != nil {
		return false
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	payload := AppError{
		App:       info.App,
		Error:     errMsg,
		Timestamp: time.Now().UnixMilli(),
	}

	data, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		return false
	}

	return r.post(r.cfg.ReportErrorURL, data)
}

// post sends a JSON POST request to the given URL. Returns true on success.
func (r *AppInfoReporter) post(url string, body []byte) bool {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := r.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}
