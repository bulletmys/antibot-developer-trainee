package middleware

import (
	"antibot-trainee/internal/limiter"
	"fmt"
	"log"
	"net"
	"net/http"
)

type RateMiddleware struct {
	RateUC limiter.UseCase
	Mask   net.IPMask
}

func (h RateMiddleware) RateLimiter(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			next.ServeHTTP(w, r)
			return
		}

		userAddr := r.Header.Get("X-Forwarded-For")

		ip := net.ParseIP(userAddr)
		subnet := ip.Mask(h.Mask)

		ok, err := h.RateUC.CheckIP(subnet)
		if err != nil {
			log.Printf("error occured while checking ip: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			maskOnes, _ := h.Mask.Size()
			log.Printf("too many requests from %v/%d", subnet, maskOnes)
			w.WriteHeader(http.StatusTooManyRequests)
			http.Error(
				w,
				fmt.Sprintf("You have exceeded the request limit, try again in %v minutes\n",
					h.RateUC.GetBlackListTTL().Minutes()),
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
