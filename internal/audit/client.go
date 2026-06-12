package audit

import (
	"arthemis-brain/internal/models"
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type WatcherClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewWatcherClient(endpoint string) *WatcherClient {
	return &WatcherClient{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *WatcherClient) Send(entry models.AuditLog) {
	payload, err := json.Marshal(entry)
	if err != nil {
		slog.Error("audit: failed to marshal log", "err", err)
		return
	}

	resp, err := c.httpClient.Post(c.endpoint+"/audit/log", "application/json", bytes.NewReader(payload))
	if err != nil {
		slog.Error("audit: failed to send log to watcher", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("audit: watcher returned non-success status", "status", resp.Status, "statusCode", resp.StatusCode)
	}
}
