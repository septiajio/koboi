package koboi

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancerService interface {
	ServeHTTP(backends []*Backend, w http.ResponseWriter, r *http.Request)
}

type LoadBalancer struct {
	balancer LoadBalancerService
	backends []*Backend
}

func New(balancer LoadBalancerService, backends []*Backend) *LoadBalancer {
	return &LoadBalancer{balancer, backends}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.balancer.ServeHTTP(lb.backends, w, r)
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
