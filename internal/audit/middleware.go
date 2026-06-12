// audit/middleware.go
package audit

import (
	"arthemis-brain/internal/models"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func AuditMiddleware(watcherClient *WatcherClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			var reqBodyStr string
			if r.Body != nil && r.Body != http.NoBody {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					if !isSensitivePath(r.URL.Path) {
						reqBodyStr = string(bodyBytes)
					}
					r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			entry := models.AuditLog{
				Timestamp:  time.Now().UTC(),
				Method:     r.Method,
				Path:       r.URL.Path,
				StatusCode: ww.Status(),
				Duration:   time.Since(start).Milliseconds(),
				UserID:     userFromContext(r.Context()),
				ReqBody:    reqBodyStr,
			}

			go watcherClient.Send(entry) // thread lightweight que o go gera
		})
	}
}

func userFromContext(ctx context.Context) string {
	id, ok := ctx.Value("user_id").(string)
	if !ok {
		return ""
	}
	return id
}

func isSensitivePath(path string) bool {
	sensitive := []string{"/login", "/register", "/password"}
	for _, s := range sensitive {
		if strings.Contains(path, s) {
			return true
		}
	}
	return false
}
