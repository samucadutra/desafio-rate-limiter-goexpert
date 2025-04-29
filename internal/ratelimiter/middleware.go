package ratelimiter

import (
	"net"
	"net/http"
	"time"
)

func Middleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			globalKey := "global_rate_limit"
			if !rl.AllowRequest(globalKey, rl.globalRateLimit, 1*time.Second) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			if rl.allowIPLimit || rl.allowTokenLimit {
				ip, _, _ := net.SplitHostPort(r.RemoteAddr)
				token := r.Header.Get("API_KEY")
				redisKey, limit, window := rl.getRedisKeyAndLimit(token, ip)

				if !rl.AllowRequest(redisKey, limit, window) {
					blockWindow := rl.getBlockWindow(token)
					rl.BlockRequest(redisKey, limit+1, blockWindow)
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
