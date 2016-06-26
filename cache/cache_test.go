package cache

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResponsePublic(t *testing.T) {

	instance, _ := aetest.NewInstance(nil)
	defer instance.Close()

	fakeR, _ := instance.NewRequest("GET", "/", nil)
	fakeW := httptest.NewRecorder()

	handler := Response(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}, 10*time.Second)

	handler(appengine.NewContext(fakeR), fakeW, fakeR)

	if fakeW.Code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, but got %d", fakeW.Code)
	}

	if fakeW.Header().Get("Cache-Control") != "public,max-age=10" {
		t.Errorf("Unexpected Cache-Control header: %s", fakeW.Header().Get("Cache-Control"))
	}

	if fakeW.Header().Get("Vary") != "" {
		t.Errorf("Unexpected Vary header: %s", fakeW.Header().Get("Vary"))
	}

	fakeCtx := appengine.NewContext(fakeR)
	fakeCtx = context.WithValue(fakeCtx, "__auth_check_ctx", "passed")

	fakeW = httptest.NewRecorder()

	handler(fakeCtx, fakeW, fakeR)

	if fakeW.Code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, but got %d", fakeW.Code)
	}

	if fakeW.Header().Get("Cache-Control") != "private,max-age=10" {
		t.Errorf("Unexpected Cache-Control header: %s", fakeW.Header().Get("Cache-Control"))
	}

	if fakeW.Header().Get("Vary") != "Authorization" {
		t.Errorf("Unexpected Vary header: %s", fakeW.Header().Get("Vary"))
	}

}
