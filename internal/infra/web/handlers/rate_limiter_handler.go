package handlers

import "net/http"

type RateLimiterHandler struct{}

func NewRateLimiterHandler() *RateLimiterHandler {
	return &RateLimiterHandler{}
}

func (h *RateLimiterHandler) HandleRateLimiterRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rate limit request\n"))
}
