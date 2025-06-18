package koboi

import (
	"log"
	"net/http"
)

type LoadBalancerService interface {
	GetNextBackend(backends []*Backend) *Backend
}

type LoadBalancer struct {
	service  LoadBalancerService
	backends []*Backend
}

func NewLoadBalancer(service LoadBalancerService, backends []*Backend) *LoadBalancer {
	return &LoadBalancer{service, backends}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.service.GetNextBackend(lb.backends)
	if backend == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Forwarding to %s", backend.URL.String())
	backend.ServeHTTP(w, r)
}

// todo: switching beetwen round-robin, weighted, and least-connection
