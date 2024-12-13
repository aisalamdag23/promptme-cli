package ratelimit

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RateLimiter struct to manage outgoing request rate
type RateLimiter struct {
	limiter *rate.Limiter
	log     *logrus.Entry
}

// NewRateLimiter initializes a rate limiter with the given rate and burst size
func NewRateLimiter(log *logrus.Entry, r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
		log:     log,
	}
}

// Do limiter wait
func (rl *RateLimiter) Do(ctx context.Context) error {
	// Wait for permission from the limiter before making the request
	if err := rl.limiter.Wait(ctx); err != nil {
		rl.log.Infoln("ratelimiter.do.wait")
		return err
	}

	rl.log.Infoln("ratelimiter.do.continue")
	return nil
}
