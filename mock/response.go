package mock

import (
	"net/http"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	return w.ResponseWriter.Write(body)
}
