package http

import (
	"net/http"
	"time"

	"github.com/mattmeyers/heimdall/logger"
)

type loggerMiddleware struct {
	logger logger.Logger
	next   http.Handler
}

type logWriter struct {
	http.ResponseWriter
	written int
	status  int
}

func (w *logWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

func newLoggingMiddleware(l logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return &loggerMiddleware{
			logger: l,
			next:   next,
		}
	}
}

func (m *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &logWriter{
		ResponseWriter: w,
		written:        0,
		status:         200,
	}

	m.next.ServeHTTP(res, r)

	if res.status < 500 {
		m.logger.Info("Response (%s): %d - %s", time.Since(start), res.status, r.URL.Path)
	} else {
		m.logger.Error("Response (%s): %d - %s", time.Since(start), res.status, r.URL.Path)
	}

}
