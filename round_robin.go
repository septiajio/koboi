package koboi

import (
	"log"
	"sync/atomic"
)

type RoundRobin struct {
	current uint64
}

// Select the next backend using round-robin
func (rr *RoundRobin) GetNextBackend(backends []*Backend) *Backend {
	index := atomic.AddUint64(&rr.current, 1)
	if len(backends) == 0 {
		log.Printf("No backend services")
		return nil
	}
	return backends[int(index)%len(backends)]
}
