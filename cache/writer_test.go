package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCachedResponse(t *testing.T) {

	w := httptest.NewRecorder()
	cacheWriter := &CachingResponseWriter{
		ResponseWriter: w,
		Public:         true,
		Duration:       100 * time.Second,
	}

	// try to mess things up
	cacheWriter.Header().Set("Cache-Control", "lol")
	cacheWriter.WriteHeader(http.StatusOK)
	cacheWriter.WriteHeader(http.StatusNotFound)
	cacheWriter.Write([]byte("OK"))
	if w.Code != http.StatusOK {
		t.Errorf("Expected response code http.StatusOK, but got %d", w.Code)
	}
	if w.Header().Get("Cache-Control") != "public,max-age=100" {
		t.Errorf("Unexpected Cache-Control header: %s", w.Header().Get("Cache-Control"))
	}
	if w.Header().Get("Pragma") != "" {
		t.Errorf("Unexpected Pragma header: %s", w.Header().Get("Pragma"))
	}
	if w.Header().Get("Vary") != "" {
		t.Errorf("Unexpected Vary header: %s", w.Header().Get("Vary"))
	}

	w = httptest.NewRecorder()
	cacheWriter = &CachingResponseWriter{
		ResponseWriter: w,
		Duration:       50 * time.Second,
	}

	cacheWriter.Write([]byte("OK"))
	if w.Code != http.StatusOK {
		t.Errorf("Expected response code http.StatusOK, but got %d", w.Code)
	}
	if w.Header().Get("Cache-Control") != "private,max-age=50" {
		t.Errorf("Unexpected Cache-Control header: %s", w.Header().Get("Cache-Control"))
	}
	if w.Header().Get("Pragma") != "" {
		t.Errorf("Unexpected Pragma header: %s", w.Header().Get("Pragma"))
	}
	if w.Header().Get("Vary") != "Authorization" {
		t.Errorf("Unexpected Vary header: %s", w.Header().Get("Vary"))
	}

}

func TestUncachedResponse(t *testing.T) {

	w := httptest.NewRecorder()
	cacheWriter := &CachingResponseWriter{
		ResponseWriter: w,
		Public:         true,
		Duration:       100 * time.Second,
	}

	// try to mess things up
	cacheWriter.Header().Set("Cache-Control", "lol")
	cacheWriter.WriteHeader(http.StatusNotFound)
	cacheWriter.Write([]byte("NotFound"))
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected response code http.StatusNotFound, but got %d", w.Code)
	}
	if w.Header().Get("Cache-Control") != "must-revalidate,max-age=0" {
		t.Errorf("Unexpected Cache-Control header: %s", w.Header().Get("Cache-Control"))
	}
	if w.Header().Get("Pragma") != "no-cache" {
		t.Errorf("Unexpected Pragma header: %s", w.Header().Get("Pragma"))
	}

}
