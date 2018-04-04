package httphelper

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// ElapsedHandler records and logs an HTTP request with the elapsed time and status code
func ElapsedHandler(logger logrus.FieldLogger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newWrappedResponseWriter(w)
		start := time.Now()
		handler.ServeHTTP(lrw, r)
		logger.Infof("%v %s %s %v", lrw.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{w, http.StatusOK}
}

func (lrw *wrappedResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
