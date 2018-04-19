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
		logMethod := logger.Infof
		if lrw.statusCode >= 400 {
			logMethod = logger.WithFields(logrus.Fields{
				"http.response": lrw.response,
			}).Errorf
		}
		logMethod("%d %s %s %s", lrw.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	response   []byte
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{w, http.StatusOK, nil}
}

func (lrw *wrappedResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *wrappedResponseWriter) Write(buff []byte) (int, error) {
	lrw.response = append(lrw.response, buff...)
	return lrw.ResponseWriter.Write(buff)
}
