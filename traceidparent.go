package traceidparent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"log"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type TraceID struct {
	next http.Handler
	name string
}

func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &TraceID{
		next: next,
		name: name,
	}, nil
}

func (m *TraceID) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	existing := req.Header.Get("Traceoriginal")
	rw.Header().Set("Traceoriginal", existing)
	actual := req.Header.Get("Traceactual")
	rw.Header().Set("Traceactual", actual)

	traceID := randomHex(16) // 16 bytes -> 32 hex chars
	spanID := randomHex(8)   // 8 bytes -> 16 hex chars

	traceparent := "00-" + traceID + "-" + spanID + "-01"

	// Add to incoming request and outgoing response
	req.Header.Set("Traceparent", traceparent)
	rw.Header().Set("Traceparent", traceparent)

	log.Printf("[TracePlugin] (existing=%s) (actual=%s) (new=%s) (path=%s)", existing, actual, traceparent, req.URL.Path)

	m.next.ServeHTTP(rw, req)
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
