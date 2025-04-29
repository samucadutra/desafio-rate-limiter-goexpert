package ratelimiter

import (
	"time"
)

type TokenConfig struct {
	Token                     string `json:"token"`
	RateLimitToken            int    `json:"rate_limit_token"`
	RateLimitWindowToken      int    `json:"rate_limit_window_token"`
	RateLimitBlockWindowToken int    `json:"rate_limit_block_window_token"`
}

type RateLimiter struct {
	storage         Storage
	limitByIP       int
	windowIP        time.Duration
	blockWindowIP   time.Duration
	tokenConfigs    map[string]TokenConfig
	globalRateLimit int
	allowIPLimit    bool
	allowTokenLimit bool
}

func NewRateLimiter(storage Storage, limitByIP int, windowIP, blockWindowIP time.Duration, tokenConfigs []TokenConfig, globalRateLimit int, allowIPLimit bool, allowTokenLimit bool) *RateLimiter {
	configMap := make(map[string]TokenConfig)
	for _, config := range tokenConfigs {
		configMap[config.Token] = config
	}

	return &RateLimiter{
		storage:         storage,
		limitByIP:       limitByIP,
		windowIP:        windowIP,
		blockWindowIP:   blockWindowIP,
		tokenConfigs:    configMap,
		globalRateLimit: globalRateLimit,
		allowIPLimit:    allowIPLimit,
		allowTokenLimit: allowTokenLimit,
	}
}

func (rl *RateLimiter) getRedisKeyAndLimit(token string, ip string) (string, int, time.Duration) {
	if rl.allowTokenLimit {
		if config, exists := rl.tokenConfigs[token]; exists {
			return token, config.RateLimitToken, time.Duration(config.RateLimitWindowToken) * time.Second
		}
	}
	if rl.allowIPLimit {
		return ip, rl.limitByIP, rl.windowIP
	}
	return "", 0, 0
}

func (rl *RateLimiter) getBlockWindow(token string) time.Duration {
	if config, exists := rl.tokenConfigs[token]; exists {
		return time.Duration(config.RateLimitBlockWindowToken) * time.Second
	}
	return rl.blockWindowIP
}

func (rl *RateLimiter) AllowRequest(key string, limit int, window time.Duration) bool {
	incr, err := rl.storage.Increment(key, window)
	if err != nil {
		return false
	}

	return incr <= int64(limit)
}

func (rl *RateLimiter) BlockRequest(key string, limit int, blockWindow time.Duration) {
	_ = rl.storage.SetValue(key, limit, blockWindow)
}
