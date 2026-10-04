package balancer

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// Balancer selects one backend for each request and proxies the request to it.
type Balancer struct {
	backends []*httputil.ReverseProxy
	next     atomic.Uint64
}

// New creates a round-robin balancer for the supplied backend URLs.
func New(targets []*url.URL) (*Balancer, error) {
	if len(targets) == 0 {
		return nil, errors.New("at least one backend is required")
	}

	backends := make([]*httputil.ReverseProxy, 0, len(targets))
	for _, target := range targets {
		if target == nil || target.Host == "" {
			return nil, errors.New("backend URL must not be nil and must have a host")
		}
		backends = append(backends, httputil.NewSingleHostReverseProxy(target))
	}
	return &Balancer{backends: backends}, nil
}

// ServeHTTP proxies the request to the next backend in round-robin order.
func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index := b.next.Add(1) - 1
	b.backends[index%uint64(len(b.backends))].ServeHTTP(w, r)
}
