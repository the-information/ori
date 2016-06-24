package rest

import (
	"net/http"
	"testing"
)

func fakeRequest(key, value string) *http.Request {
	return &http.Request{
		Header: map[string][]string{
			key: []string{value},
		},
	}
}

func TestAcceptsJson(t *testing.T) {

	if !AcceptsJson(fakeRequest("", "")) {
		t.Errorf("AcceptsJson should return true for an empty Accept, but got false")
	}

	if AcceptsJson(fakeRequest("Accept", "text/plain")) {
		t.Errorf("AcceptsJson should return false for Accept: text/plain, but got true")
	}

	if !AcceptsJson(fakeRequest("Accept", "application/json")) {
		t.Errorf("AcceptsJson should return true for Accept: application/json, but got false")
	}

	if !AcceptsJson(fakeRequest("Accept", "*/*")) {
		t.Errorf("AcceptsJson should return true for Accept: */*, but got false")
	}

	// Accept headers should NEVER use ;charset=UTF-8, but should
	// set the Accept-Charset header instead.
	if AcceptsJson(fakeRequest("Accept", "application/json;charset=utf-8")) {
		t.Errorf("AcceptsJson should return false for Accept: application/json; charset=UTF-8, but got true")
	}

}

func TestAcceptsUtf8(t *testing.T) {

	if !AcceptsUtf8(fakeRequest("", "")) {
		t.Errorf("AcceptsUtf8 should return true for empty Accept-Charset, but got false")
	}

	if !AcceptsUtf8(fakeRequest("Accept-Charset", "*")) {
		t.Errorf("AcceptsUtf8 should return true for Accept-Charset: *, but got false")
	}

	if AcceptsUtf8(fakeRequest("Accept-Charset", "ISO-8859-1")) {
		t.Errorf("AcceptsUtf8 should return false for Accept-Charset: UTF-8, but got true")
	}

	if !AcceptsUtf8(fakeRequest("Accept-Charset", "UTF-8")) {
		t.Errorf("AcceptsUtf8 should return true for Accept-Charset: UTF-8, but got false")
	}

}

func TestContentIsJson(t *testing.T) {

	if !ContentIsJson(fakeRequest("Content-Type", "application/json")) {
		t.Errorf("ContentIsJson should return true for Content-Type: application/json, but got false")
	}

	if !ContentIsJson(fakeRequest("Content-Type", "application/json")) {
		t.Errorf("ContentIsJson should return true for Content-Type: application/json, but got false")
	}

	if ContentIsJson(fakeRequest("Content-Type", "application/x-www-form-urlencoded")) {
		t.Errorf("ContentIsJson should return false for Content-Type: application/json;charset=ISO-8859-1, but got false")
	}

	if ContentIsJson(fakeRequest("Content-Type", "application/json;charset=ISO-8859-1")) {
		t.Errorf("ContentIsJson should return false for Content-Type: application/json;charset=ISO-8859-1, but got false")
	}

}

func TestHasValidOrigin(t *testing.T) {

	if !HasValidOrigin(fakeRequest("Origin", "http://foo.bar.com"), "bar.com") {
		t.Errorf("HasValidOrigin should return true with Origin: http://foo.bar.com and domain suffix bar.com, but got false")
	}

	if HasValidOrigin(fakeRequest("Origin", "http://foo.bar.com"), "quux.com") {
		t.Errorf("HasValidOrigin should return false with Origin: http://foo.bar.com and domain suffix quux.com, but got true")
	}

	if !HasValidOrigin(fakeRequest("", ""), "quux.com") {
		t.Errorf("HasValidOrigin should return true for an empty Origin, but got false")
	}

}
