package cache

import (
	"fmt"
	"net/http"
	"time"
)

type CachingResponseWriter struct {
	http.ResponseWriter
	Duration   time.Duration
	Public     bool
	wasWritten bool
}

func (w *CachingResponseWriter) WriteHeader(code int) {

	if w.wasWritten {
		return
	}

	w.wasWritten = true

	var cacheLevel string
	if w.Public {
		cacheLevel = "public"
	} else {
		cacheLevel = "private"
	}

	// set the cache duration headers
	// first thing to note: if this isn't a 200 OK, don't cache it!
	if code != http.StatusOK {

		w.Header().Set("Cache-Control", "must-revalidate,max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.ResponseWriter.WriteHeader(code)

	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("%s,max-age=%d", cacheLevel, w.Duration/time.Second))

		if !w.Public {
			w.Header().Set("Vary", "Authorization")
		}

		w.ResponseWriter.WriteHeader(code)
	}

}

func (w *CachingResponseWriter) Write(b []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	return w.ResponseWriter.Write(b)
}
