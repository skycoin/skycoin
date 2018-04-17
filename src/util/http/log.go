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
		if lrw.statusCode == 200 || lrw.statusCode == 304 {
			logger.Infof("%v %s %s %v", lrw.statusCode, r.Method, r.URL.Path, time.Since(start))
		} else {
			if lrw.statusCode == 400 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Bad Request", r.Method, r.URL.Path, time.Since(start))
			}
			if lrw.statusCode == 403 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Forbidden", r.Method, r.URL.Path, time.Since(start))
			}
			if lrw.statusCode == 404 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Not Found", r.Method, r.URL.Path, time.Since(start))
			}
			if lrw.statusCode == 405 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Method Not Allowed", r.Method, r.URL.Path, time.Since(start))
			}
			if lrw.statusCode == 500 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Internal Server Error", r.Method, r.URL.Path, time.Since(start))
			}
			if lrw.statusCode == 501 {
				logger.Errorf("%v %s %s %s %v", lrw.statusCode, "Not Implemented", r.Method, r.URL.Path, time.Since(start))
			}

		}

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
