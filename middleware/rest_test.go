package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
func requestWithHeader(method string, h http.Header) *http.Request {
	r, _ := http.NewRequest(method, "http://example.com", nil)
	r.Header = h
	return r
}

func TestRest(t *testing.T) {

	testHandler := RestJsonValidateHandler{
		ValidOriginSuffix: "example.com",
	}

	w := httptest.NewRecorder()

	// bad accept header should get http.NotAcceptable
	testHandler.ServeHTTP(w, requestWithHeader("POST", map[string][]string{
		"Accept": []string{"application/xml"},
	}))

	if w.Code != http.StatusNotAcceptable {
		t.Errorf("Wrong status code, wanted %d, got %d", http.StatusNotAcceptable, w.Code)
	}

	w = httptest.NewRecorder()

	// bad content-type should get http.UnsupportedMediaType
	testHandler.ServeHTTP(w, requestWithHeader("POST", map[string][]string{
		"Content-Type": []string{"application/xml"},
	}))

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Wrong status code, wanted %d, got %d", http.StatusUnsupportedMediaType, w.Code)
	}

	w = httptest.NewRecorder()

	// bad content-type should be OK with a GET or HEAD
	testHandler.ServeHTTP(w, requestWithHeader("HEAD", map[string][]string{
		"Content-Type": []string{"application/xml"},
	}))

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code, wanted %d, got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()

	// bad Origin should get http.BadRequest
	testHandler.ServeHTTP(w, requestWithHeader("POST", map[string][]string{
		"Content-Type": []string{"application/json"},
		"Origin":       []string{"http://baddomain.com"},
	}))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Wrong status code, wanted %d, got %d", http.StatusBadRequest, w.Code)
	}

	w = httptest.NewRecorder()

	// good Origin should result in Access-Control-Allow-Origin being set on the response
	testHandler.ServeHTTP(w, requestWithHeader("POST", map[string][]string{
		"Content-Type": []string{"application/json"},
		"Origin":       []string{"http://sub.example.com"},
	}))

	if w.Header().Get("Access-Control-Allow-Origin") != "http://sub.example.com" {
		t.Errorf("Access-Control-Allow-Origin was not set when it should have been")
	}

	w = httptest.NewRecorder()

	// no Origin should result in no Access-Control-Allow-Origin on the response
	testHandler.ServeHTTP(w, requestWithHeader("POST", map[string][]string{
		"Content-Type": []string{"application/json"},
	}))

	if w.Header().Get("X-REST-OK") != "OK" {
		t.Errorf("expected X-REST-OK to be OK, but got %s", w.Header().Get("X-REST-OK"))
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Access-Control-Allow-Origin was set when it should not have been")
	}
}
*/
