package qos

import (
	"sync"

	"golang.org/x/time/rate"
)

type IDRateLimiter struct {
	ids   map[uint]*rate.Limiter
	mu    *sync.RWMutex
	rate  rate.Limit // per second
	burst int
}

func NewIDRateLimiter(r rate.Limit, b int) *IDRateLimiter {
	i := &IDRateLimiter{
		ids:   make(map[uint]*rate.Limiter),
		mu:    &sync.RWMutex{},
		rate:  r,
		burst: b,
	}

	return i
}

func (i *IDRateLimiter) AddID(id uint) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ids[id] = limiter

	return limiter
}

func (i *IDRateLimiter) GetLimiter(id uint) *rate.Limiter {
	i.mu.RLock()
	defer i.mu.RUnlock()

	limiter, exist := i.ids[id]
	if exist {
		return limiter
	}

	return i.AddID(id)
}
