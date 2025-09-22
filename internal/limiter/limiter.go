package limiter

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	rate     int
	interval time.Duration
	tokens   int
	mu       sync.Mutex
	lastTime time.Time
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		tokens:   rate,
		lastTime: time.Now(),
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime)
	
	tokensToAdd := int(elapsed / rl.interval)
	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.rate)
		rl.lastTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return nil
	}

	waitTime := rl.interval - elapsed
	if waitTime > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			rl.tokens = rl.rate - 1
			rl.lastTime = time.Now()
		}
	}

	return nil
}

func (rl *RateLimiter) SetRate(rate int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
}

func (rl *RateLimiter) GetRate() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.rate
}

type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
	Backoff    BackoffStrategy
}

type BackoffStrategy interface {
	GetDelay(attempt int) time.Duration
}

type LinearBackoff struct {
	BaseDelay time.Duration
}

func (lb *LinearBackoff) GetDelay(attempt int) time.Duration {
	return lb.BaseDelay * time.Duration(attempt)
}

type ExponentialBackoff struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
}

func (eb *ExponentialBackoff) GetDelay(attempt int) time.Duration {
	delay := eb.BaseDelay * time.Duration(1<<uint(attempt-1))
	if delay > eb.MaxDelay {
		delay = eb.MaxDelay
	}
	return delay
}

type Retryer struct {
	config RetryConfig
}

func NewRetryer(config RetryConfig) *Retryer {
	return &Retryer{config: config}
}

func (r *Retryer) Execute(ctx context.Context, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := r.config.Backoff.GetDelay(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return lastErr
}

func (r *Retryer) ExecuteWithResult[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	var result T
	var lastErr error
	
	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := r.config.Backoff.GetDelay(attempt)
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(delay):
			}
		}

		res, err := fn()
		if err == nil {
			return res, nil
		}

		result = res
		lastErr = err
	}

	return result, lastErr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
