package httpapi

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type (
	logMW struct {
		next http.Handler // loggerMW is the first handler and then it calls servehttp of the next middleware or final handler
		log  *log.Logger  // Actual logger
	}

	logWriter struct {
		w   http.ResponseWriter
		hdr int
		buf []byte
	}
)

// Ctor
func NewLogMW(next http.Handler) *logMW {
	return &logMW{
		next: next,
		log:  log.Default(),
	}
}

func newLogWriter(w http.ResponseWriter) *logWriter {
	return &logWriter{
		w:   w,
		hdr: http.StatusOK,
		buf: make([]byte, 0, 1024),
	}
}

// logWriter must implement ResponseWriter to be able to save response data
// Header
func (lw *logWriter) Header() http.Header {
	return lw.w.Header()
}

// Write
func (lw *logWriter) Write(in []byte) (int, error) {
	lw.buf = append(lw.buf, in...)
	return lw.w.Write(in)
}

// WriteHeader
func (lw *logWriter) WriteHeader(statusCode int) {
	lw.hdr = statusCode
	lw.w.WriteHeader(statusCode)
}

// To be fed to server instead of some handler => logMW must implement ServeHTTP
func (mw *logMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Dumping incoming request to []byte var
	reqDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		mw.log.Printf("ERROR: dump request: %v\n", err)
	} else {
		mw.log.Printf("REQUEST:\n%s\n", string(reqDump))
	}
	// Initializing response log writer with base response writer
	lw := newLogWriter(w)
	// Logging response by calling next handler and writing data to log writer (put instead of original responseWriter w)
	mw.next.ServeHTTP(lw, r)
	// Print response data saved to
	mw.log.Printf("RESPONSE:\nSTATUS:%d\nBODY:%s\n", lw.hdr, string(lw.buf))
}
