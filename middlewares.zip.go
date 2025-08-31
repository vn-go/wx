package wx

import (
	"compress/gzip"
	"net/http"
	"strings"
)

var zip = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// Client không hỗ trợ gzip
		next.ServeHTTP(w, r)
		return
	}

	// Gửi header báo là đã nén gzip
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	gzrw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	next.ServeHTTP(gzrw, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
