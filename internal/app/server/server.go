package server

import (
	rateLimiter "antibot-trainee/internal/limiter/delivery/http"
	"antibot-trainee/internal/limiter/repository"
	"antibot-trainee/internal/limiter/usecase"
	"antibot-trainee/internal/middleware"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func RunServer() {
	maskSize, err := strconv.Atoi(os.Getenv("MASK_SIZE"))
	if err != nil {
		log.Fatalf("failed to parse mask from config: %v", err)
	}

	ipv4Mask := net.CIDRMask(maskSize, 32)

	blTTL, err := strconv.Atoi(os.Getenv("BLACKLIST_TTL"))
	if err != nil {
		log.Fatalf("failed to parse black list ttl from config: %v", err)
	}

	rTTL, err := strconv.Atoi(os.Getenv("REQUEST_TTL"))
	if err != nil {
		log.Fatalf("failed to parse request ttl from config: %v", err)
	}

	blTTLDuration := time.Duration(blTTL)
	rTTLDuration := time.Duration(rTTL)

	rUC := &usecase.RateUseCase{
		RateRepo: repository.NewMapRepo(100, time.Second*blTTLDuration, time.Second*rTTLDuration),
	}

	h := middleware.RateMiddleware{
		RateUC: rUC,
		Mask:   ipv4Mask,
	}

	handler := rateLimiter.RateHandler{RateUC: rUC}
	http.Handle("/data", h.RateLimiter(handler.Handler))

	httpAddr := os.Getenv("HTTP_ADDR")
	log.Printf("starting server at %v", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
