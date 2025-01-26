package utils

import (
    "context"
    "time"
    "golang.org/x/time/rate"
)

type RateLimiterConfig struct {
    Limit     rate.Limit
    BurstSize int
}

func NewRateLimiter(limit rate.Limit, burstSize int) *rate.Limiter {
    return rate.NewLimiter(limit, burstSize)
}

func WaitForRateLimit(ctx context.Context, limiter *rate.Limiter) error {
    return limiter.Wait(ctx)
}

// Predefined configurations
var (
    LenientRateLimiter = RateLimiterConfig{
        Limit:     rate.Every(time.Minute / 10),  // 10 requests per minute
        BurstSize: 10,
    }

    StrictRateLimiter = RateLimiterConfig{
        Limit:     rate.Every(10 * time.Second), // Fewer requests, more conservative
        BurstSize: 5,
    }
)
// package utils

// import (
//     "context"
//     "time"
//     "golang.org/x/time/rate"
// )

// type RateLimiterConfig struct {
//     Limit     rate.Limit
//     BurstSize int
// }

// func NewRateLimiter(limit rate.Limit, burstSize int) *rate.Limiter {
//     return rate.NewLimiter(limit, burstSize)
// }

// func WaitForRateLimit(ctx context.Context, limiter *rate.Limiter) error {
//     return limiter.Wait(ctx)
// }

// // Predefined configurations
// var (
//     LenientRateLimiter = RateLimiterConfig{
//         Limit:     rate.Every(6 * time.Second),
//         BurstSize: 10,
//     }

//     StrictRateLimiter = RateLimiterConfig{
//         Limit:     rate.Every(10 * time.Second),
//         BurstSize: 5,
//     }
// )