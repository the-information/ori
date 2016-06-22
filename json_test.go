package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {

	x := struct {
		Message string `json:"message"`
	}{}

	r, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(`{"message": "hello"}`))

	if err := ReadJSON(r, &x); err != nil {
		t.Errorf("Got unexpected error %s", err)
	} else if x.Message != "hello" {
		t.Errorf("Got unexpected message %s", x.Message)
	}

}

func TestWriteJSON(t *testing.T) {

	x := struct {
		Quantity int `json:"quantity"`
	}{
		5,
	}

	w := httptest.NewRecorder()
	WriteJSON(w, &x)
	if strings.TrimSpace(w.Body.String()) != `{"quantity":5}` {
		t.Errorf("Got unexpected body `%s`", w.Body)
	}

	w = httptest.NewRecorder()
	WriteJSON(w, &ErrNotFound)
	if w.Code != http.StatusNotFound {
		t.Errorf("Got unexpected response code, wanted 404, got %d", w.Code)
	} else if strings.TrimSpace(w.Body.String()) != `{"message":"The requested resource could not be located."}` {
		t.Errorf("Got unexpected response %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	WriteJSON(w, errors.New("Wat"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Got unexpected response code, wanted 500, got %d", w.Code)
	} else if strings.TrimSpace(w.Body.String()) != `{"message":"Wat"}` {

	}

}
