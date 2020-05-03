package http

import (
	"antibot-trainee/internal/limiter/usecase"
	"fmt"
	"log"
	"net/http"
)

type RateHandler struct {
	RateUC *usecase.RateUseCase
}

func (h *RateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		keys, ok := r.URL.Query()["prefix"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'prefix' is missing")
			return
		}
		h.RateUC.ResetLimit(keys[0])
	default:
		if _, err := fmt.Fprint(w, "OK\n"); err != nil {
			log.Printf("failed to send response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
