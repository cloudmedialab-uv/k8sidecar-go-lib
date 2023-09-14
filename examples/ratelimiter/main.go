package ratelimiter

import (
	"net/http"
	"sync"
	"time"

	sidecar "github.com/k8sidecar/go-lib"
)

type RateLimiter struct {
	visitors       map[string]*Visitor
	rate           int
	refreshTimeout time.Duration
	mu             sync.Mutex
}

type Visitor struct {
	lastAccessed time.Time
	requestCount int
}

func NewRateLimiter(rate int, refreshTimeout time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors:       make(map[string]*Visitor),
		rate:           rate,
		refreshTimeout: refreshTimeout,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists || time.Now().Sub(v.lastAccessed) > rl.refreshTimeout {
		rl.visitors[ip] = &Visitor{lastAccessed: time.Now(), requestCount: 1}
		return true
	}

	if v.requestCount < rl.rate {
		v.requestCount++
		v.lastAccessed = time.Now()
		return true
	}

	return false
}

func (rl *RateLimiter) RateLimiterMiddleware(req *http.Request, res http.ResponseWriter, chain *sidecar.FilterChain) {
	clientIP := req.RemoteAddr
	if !rl.Allow(clientIP) {
		res.WriteHeader(http.StatusTooManyRequests)
		res.Write([]byte("Too many requests"))
		return
	}
	chain.Next()
}

func main() {
	rl := NewRateLimiter(100, 10*time.Minute)

	filter := &sidecar.SidecarFilter{
		TriFunction: rl.RateLimiterMiddleware,
	}
	filter.Listen()

}
