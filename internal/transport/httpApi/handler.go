package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func (h *Handler) CalcHandler(w http.ResponseWriter, r *http.Request) {
	req := new(domain.JsonReq)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		http.Error(w, pkg.ErrIvalJson.Error(), http.StatusBadRequest)
		log.Println(pkg.ErrIvalJson.Error()+":", req, "error:", err)
		return
	}
	log.Println("expression received:", req)

	result, err := pkg.Calc(req.Expression)

	if err != nil {

		if errors.Is(err, pkg.ErrInvalidExpr) {
			log.Println("Invalid expression on the calculating")
			w.WriteHeader(http.StatusUnprocessableEntity)
			var vid = map[string]string{"error": err.Error()}

			_ = json.NewEncoder(w).Encode(vid)

			return

		}
		log.Println("error: Internal server error", err)

		w.WriteHeader(http.StatusInternalServerError)
		var vid = map[string]string{"error": "Internal server error"}

		_ = json.NewEncoder(w).Encode(vid)
		return

	}
	w.WriteHeader(http.StatusOK)
	var vid = map[string]string{"result": strconv.Itoa(int(result))}

	_ = json.NewEncoder(w).Encode(vid)

	log.Printf("result: %f", result)

}

// Мидлвaря няшка
func LoggingMiddleware(DefaLogger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			DefaLogger.Info("finished", slog.Group("request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
				slog.Int("status", http.StatusOK),
				slog.Duration("duration", time.Second),
			)
			next.ServeHTTP(w, r)

		})
	}
}

func NewHandler() *Handler {
	return &Handler{}
}
