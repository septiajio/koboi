package koboi

import (
	"log"
	"net/http"
	"sync/atomic"
)

type RoundRobin struct {
	current uint64
	config  *Config
}

// Select the next backend using round-robin
func (lb *RoundRobin) getNextBackend(backends []*Backend) *Backend {
	index := atomic.AddUint64(&lb.current, 1)
	if len(backends) == 0 {
		log.Printf("No backend services")
		return nil
	}
	return backends[int(index)%len(backends)]
}

// Handle and forward the request
func (lb *RoundRobin) ServeHTTP(backends []*Backend, w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend(backends)
	log.Printf("Forwarding request to: %s", backend.URL.String())
	backend.ReverseProxy.ServeHTTP(w, r)
}
