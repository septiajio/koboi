package koboi

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	current  uint64
}

// Select the next backend using round-robin
func (lb *LoadBalancer) getNextBackend() *Backend {
	index := atomic.AddUint64(&lb.current, 1)
	return lb.backends[int(index)%len(lb.backends)]
}

// Handle and forward the request
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	log.Printf("Forwarding request to: %s", backend.URL.String())
	backend.ReverseProxy.ServeHTTP(w, r)
}

func NewBackend(rawurl string) *Backend {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		log.Fatalf("Invalid backend URL: %v", err)
	}
	return &Backend{
		URL:          parsedURL,
		Alive:        true,
		ReverseProxy: newServeMuxProxy(parsedURL),
	}
}

// Create reverse proxy for backend
func newServeMuxProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = target.Path
			req.Host = target.Host
		},
	}
}

// todo: switching beetwen round-robin, weighted, and least-connection
