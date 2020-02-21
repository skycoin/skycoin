package httphelper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWithError struct {
	Data  *json.RawMessage `json:"data,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

// ElapsedHandler records and logs an HTTP request with the elapsed time and status code
func ElapsedHandler(logger logrus.FieldLogger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newWrappedResponseWriter(w)
		start := time.Now()
		handler.ServeHTTP(lrw, r)
		logMethod := logger.Infof
		if lrw.statusCode >= 400 {
			var errMsg string
			if strings.Contains(r.URL.RequestURI(), "v2") {
				rsp := responseWithError{}
				err := json.NewDecoder(strings.NewReader(lrw.response.String())).Decode(&rsp)
				// Incorrect URI address would return "404 Not Found" error, which would fail
				// the json decoding.
				if err != nil && strings.Contains(err.Error(), "json: cannot unmarshal") {
					errMsg = lrw.response.String()
				} else {
					errMsg = rsp.Error.Message
				}
			}

			logMethod = logger.WithFields(logrus.Fields{
				"body": strings.TrimSpace(errMsg),
			}).Errorf
		}
		logMethod("%d %s %s %s", lrw.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	response   bytes.Buffer
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		response:       bytes.Buffer{},
	}
}

func (lrw *wrappedResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *wrappedResponseWriter) Write(buff []byte) (int, error) {
	retVal, err := lrw.ResponseWriter.Write(buff)
	if lrw.statusCode >= 400 {
		lrw.response.Write(buff)
	}
	return retVal, err
}
