package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// NewHTTPHandler creates a video server for videoPath.
//
// The returned handler exposes:
//
//	GET /video   - serves the video, including HTTP byte-range requests
//	HEAD /video  - returns the video metadata without its body
//	GET /healthz - reports whether the server is running
func NewHTTPHandler(videoPath string) (http.Handler, error) {
	info, err := os.Stat(videoPath)
	if err != nil {
		return nil, fmt.Errorf("inspect video %q: %w", videoPath, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("video path %q is not a regular file", videoPath)
	}

	absVideoPath, err := filepath.Abs(videoPath)
	if err != nil {
		return nil, fmt.Errorf("resolve video path %q: %w", videoPath, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /video", serveVideo(absVideoPath))
	mux.HandleFunc("HEAD /video", serveVideo(absVideoPath))

	return mux, nil
}

func serveVideo(videoPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ServeFile uses ServeContent internally, so Range and If-Modified-Since
		// requests work without loading the complete video into memory.
		http.ServeFile(w, r, videoPath)
	}
}

// RunHTTP validates the configured video and starts the HTTP server.
func RunHTTP(addr, videoPath string) error {
	handler, err := NewHTTPHandler(videoPath)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("serving %s at http://%s/video", videoPath, addr)
	err = httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
