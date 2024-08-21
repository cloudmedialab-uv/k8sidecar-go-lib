package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	sidecar "github.com/cloudmedialab-uv/k8sidecar-go-lib"
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
	if !exists || time.Since(v.lastAccessed) > rl.refreshTimeout {
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

func getClientIP(req *http.Request) string {
	forwarded := req.Header.Get("X-Forwarded-For")
	if forwarded != "" {

		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}

	return ip
}

func (rl *RateLimiter) RateLimiterMiddleware(req *http.Request, res http.ResponseWriter, chain *sidecar.FilterChain) {
	clientIP := getClientIP(req)
	log.Println("client IP " + clientIP)
	if !rl.Allow(clientIP) {
		res.WriteHeader(http.StatusTooManyRequests)
		res.Write([]byte("Too many requests"))
		return
	}
	chain.Next()
}

func main() {

	rateString, exist := os.LookupEnv("RATE")

	if !exist {
		log.Println("RATE env: using default 100")
		rateString = "100"
	}

	rate, err := strconv.Atoi(rateString)

	if err != nil {
		log.Fatal("Error converting RATE env: using default 100")
		rate = 100
	}

        log.Println("RATE (max requests per second): " + strconv.Itoa(rate))

	rl := NewRateLimiter(rate, time.Second)

	filter := &sidecar.SidecarFilter{
		TriFunction: rl.RateLimiterMiddleware,
	}
	filter.Listen()
}
