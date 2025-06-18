package koboi

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL             *url.URL
	Weight          int // For weighted strategy
	CurrentWeight   int
	EffectiveWeight int
	Alive           bool
	ReverseProxy    *httputil.ReverseProxy
	mux             sync.RWMutex
}

func NewBackend(rawURL string, weight int) *Backend {
	parsedURL, _ := url.Parse(rawURL)
	return &Backend{
		URL:             parsedURL,
		Weight:          weight,
		EffectiveWeight: weight,
		Alive:           true,
		ReverseProxy:    httputil.NewSingleHostReverseProxy(parsedURL),
	}
}

func (b *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.ReverseProxy.ServeHTTP(w, r)
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}
