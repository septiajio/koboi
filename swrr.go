package koboi

import "sync"

type Weighted struct {
	mux sync.Mutex
}

// Smooth Weighted Round Robin
func (lb *Weighted) GetNextBackend(backends []*Backend) *Backend {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	var best *Backend
	total := 0

	for _, b := range backends {
		if !b.IsAlive() {
			continue
		}
		b.CurrentWeight += b.EffectiveWeight
		total += b.EffectiveWeight

		if best == nil || b.CurrentWeight > best.CurrentWeight {
			best = b
		}
	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
