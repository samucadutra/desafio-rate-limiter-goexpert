package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/samucadutra/desafio-rate-limiter-goexpert/configs"
	"github.com/samucadutra/desafio-rate-limiter-goexpert/internal/infra/web/handlers"
	"github.com/samucadutra/desafio-rate-limiter-goexpert/internal/infra/web/webserver"
	"github.com/samucadutra/desafio-rate-limiter-goexpert/internal/ratelimiter"
	"strconv"
	"time"
)

func main() {
	appConfig, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: appConfig.RedisHost + ":" + appConfig.RedisPort,
	})
	defer redisClient.Close()

	storage := ratelimiter.NewRedisStorage(redisClient)

	rateLimitWindowIP, err := strconv.Atoi(appConfig.RateLimitWindowIP)
	if err != nil {
		panic(err)
	}
	rateLimitIP, err := strconv.Atoi(appConfig.RateLimitIp)
	if err != nil {
		panic(err)
	}

	rateLimitBlockWindowIP, err := strconv.Atoi(appConfig.RateLimitBlockWindowIP)
	if err != nil {
		panic(err)
	}

	globalRateLimit, err := strconv.Atoi(appConfig.GlobalRateLimit)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse GLOBAL_RATE_LIMIT: %v", err))
	}

	var tokenConfigs []ratelimiter.TokenConfig
	tokensConfigLimit := appConfig.TokensConfigLimit
	if err := json.Unmarshal([]byte(tokensConfigLimit), &tokenConfigs); err != nil {
		panic(fmt.Sprintf("Failed to parse TOKENS_CONFIG_LIMIT: %v", err))
	}

	allowIPLimit, err := strconv.ParseBool(appConfig.AllowIPLimit)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse ALLOW_IP_LIMIT: %v", err))
	}

	allowTokenLimit, err := strconv.ParseBool(appConfig.AllowTokenLimit)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse ALLOW_TOKEN_LIMIT: %v", err))
	}

	rateLimiter := ratelimiter.NewRateLimiter(
		storage,
		rateLimitIP,
		time.Duration(rateLimitWindowIP)*time.Second,
		time.Duration(rateLimitBlockWindowIP)*time.Second,
		tokenConfigs,
		globalRateLimit,
		allowIPLimit,
		allowTokenLimit,
	)

	rateLimiterHandler := handlers.NewRateLimiterHandler()

	webServer := webserver.NewWebServer(appConfig.WebServerPort)

	webServer.AddCustomizedMiddleware(ratelimiter.Middleware(rateLimiter))

	webServer.AddHandler("/", rateLimiterHandler.HandleRateLimiterRequest)
	fmt.Println("Starting web server on port", appConfig.WebServerPort)
	webServer.Start()
}
