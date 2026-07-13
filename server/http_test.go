package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testVideo = "0123456789video-data"

func TestNewHTTPHandlerRejectsMissingVideo(t *testing.T) {
	_, err := NewHTTPHandler(filepath.Join(t.TempDir(), "missing.mp4"))
	if err == nil {
		t.Fatal("NewHTTPHandler() error = nil, want an error")
	}
}

func TestHTTPHandler(t *testing.T) {
	handler := newTestHTTPHandler(t)

	tests := []struct {
		name       string
		method     string
		path       string
		rangeValue string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "health check",
			method:     http.MethodGet,
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   "ok\n",
		},
		{
			name:       "complete video",
			method:     http.MethodGet,
			path:       "/video",
			wantStatus: http.StatusOK,
			wantBody:   testVideo,
		},
		{
			name:       "video byte range",
			method:     http.MethodGet,
			path:       "/video",
			rangeValue: "bytes=2-5",
			wantStatus: http.StatusPartialContent,
			wantBody:   "2345",
		},
		{
			name:       "video metadata",
			method:     http.MethodHead,
			path:       "/video",
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "unknown route",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
		{
			name:       "unsupported method",
			method:     http.MethodPost,
			path:       "/video",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method Not Allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.rangeValue != "" {
				req.Header.Set("Range", tt.rangeValue)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, req)

			result := response.Result()
			defer result.Body.Close()
			body, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("read response: %v", err)
			}

			if result.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", result.StatusCode, tt.wantStatus)
			}
			if string(body) != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}

			if tt.rangeValue != "" {
				if got := result.Header.Get("Content-Range"); got != "bytes 2-5/20" {
					t.Errorf("Content-Range = %q, want %q", got, "bytes 2-5/20")
				}
			}
		})
	}
}

func TestVideoIsServedWithoutLoadingItAtStartup(t *testing.T) {
	path := writeTestVideo(t)
	handler, err := NewHTTPHandler(path)
	if err != nil {
		t.Fatalf("NewHTTPHandler() error = %v", err)
	}

	updatedVideo := strings.Repeat("x", len(testVideo))
	if err := os.WriteFile(path, []byte(updatedVideo), 0o600); err != nil {
		t.Fatalf("update video: %v", err)
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/video", nil))

	if response.Body.String() != updatedVideo {
		t.Fatalf("body = %q, want updated file %q", response.Body.String(), updatedVideo)
	}
}

func newTestHTTPHandler(t *testing.T) http.Handler {
	t.Helper()
	handler, err := NewHTTPHandler(writeTestVideo(t))
	if err != nil {
		t.Fatalf("NewHTTPHandler() error = %v", err)
	}
	return handler
}

func writeTestVideo(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.mp4")
	if err := os.WriteFile(path, []byte(testVideo), 0o600); err != nil {
		t.Fatalf("write test video: %v", err)
	}
	return path
}
