package server

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	// ErrTooManyRequests is returned when the rate limit is exceeded.
	ErrTooManyRequests = errors.New("too many requests")
)

const (
	// defaultRateLimiterRate is the default rate limiter rate.
	defaultRateLimiterRate = 1
	// defaultRateLimiterBurst is the default rate limiter burst.
	defaultRateLimiterBurst = 3
	// defaultRateLimiterTTL is the default rate limiter time-to-live.
	defaultRateLimiterTTL = 5 * time.Minute
	// defaultRateLimiterCleanupInterval is the default rate limiter cleanup interval.
	defaultRateLimiterCleanupInterval = 10 * time.Second
)

// rateLimiter represents a rate limiter for an IP address.
type rateLimiter struct {
	limiter *rate.Limiter
	created time.Time
}

// rateLimiters contains a map of rate limiters for each IP address,
// and the rate limit options.
type rateLimiters struct {
	limiters        map[string]*rateLimiter
	rate            rate.Limit
	burst           int
	ttl             time.Duration
	cleanupInterval time.Duration
	stop            chan struct{}
	mu              sync.Mutex
}

// get a rate limiter for the given IP address. If a rate limiter does not exist for the IP address,
// a new one is created.
func (c *rateLimiters) get(ip string) *rateLimiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	rl, ok := c.limiters[ip]
	if !ok {
		c.limiters[ip] = &rateLimiter{
			limiter: rate.NewLimiter(c.rate, c.burst),
			created: time.Now(),
		}
		return c.limiters[ip]
	}
	return rl
}

// close the rate limiters and stop the cleanup goroutine.
func (c *rateLimiters) close() error {
	c.stop <- struct{}{}
	close(c.stop)
	return nil
}

// cleanup removes rate limiters that have expired.
func (c *rateLimiters) cleanup() {
	for {
		select {
		case <-time.After(c.cleanupInterval):
			c.mu.Lock()
			for ip, limiter := range c.limiters {
				if time.Since(limiter.created) > c.ttl {
					delete(c.limiters, ip)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

// rateLimiterOptions contains the options for the rate limiter middleware.
type rateLimiterOptions struct {
	rate            rate.Limit
	burst           int
	ttl             time.Duration
	cleanupInterval time.Duration
}

// rateLimiterOption is a function that configures the rate limiter options.
type rateLimiterOption func(o *rateLimiterOptions)

// rateLimitHandler is a middleware that limits the number of requests that can be made to the server
// on a per-IP basis.
func rateLimitHandler(options ...rateLimiterOption) func(next http.Handler) (http.Handler, func() error) {
	opts := rateLimiterOptions{
		rate:            defaultRateLimiterRate,
		burst:           defaultRateLimiterBurst,
		ttl:             defaultRateLimiterTTL,
		cleanupInterval: defaultRateLimiterCleanupInterval,
	}
	for _, option := range options {
		option(&opts)
	}

	rateLimiters := &rateLimiters{
		limiters:        make(map[string]*rateLimiter),
		rate:            opts.rate,
		burst:           opts.burst,
		ttl:             opts.ttl,
		cleanupInterval: opts.cleanupInterval,
		stop:            make(chan struct{}),
		mu:              sync.Mutex{},
	}

	go rateLimiters.cleanup()

	return func(next http.Handler) (http.Handler, func() error) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rl := rateLimiters.get(resolveIP(r))
			if !rl.limiter.Allow() {
				writeError(w, http.StatusTooManyRequests, ErrTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		}), rateLimiters.close
	}
}

// withRateLimiterRate sets the rate limiter rate.
func withRateLimiterRate(r float64) rateLimiterOption {
	return func(o *rateLimiterOptions) {
		if r != 0 {
			o.rate = rate.Limit(r)
		}
	}
}

// withRateLimiterBurst sets the rate limiter burst.
func withRateLimiterBurst(burst int) rateLimiterOption {
	return func(o *rateLimiterOptions) {
		if burst != 0 {
			o.burst = burst
		}
	}
}

// withRateLimiterTTL sets the rate limiter time-to-live.
func withRateLimiterTTL(ttl time.Duration) rateLimiterOption {
	return func(o *rateLimiterOptions) {
		if ttl != 0 {
			o.ttl = ttl
		}
	}
}

// withRateLimiterCleanupInterval sets the rate limiter cleanup interval.
func withRateLimiterCleanupInterval(interval time.Duration) rateLimiterOption {
	return func(o *rateLimiterOptions) {
		if interval != 0 {
			o.cleanupInterval = interval
		}
	}
}
