package httpapi

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	"github.com/gorilla/mux"
)

type key string

const Key = key("Logger")

type Router struct {
	Router  *mux.Router
	Handler Handler

	config *config.Config
}

func NewRouter(ctx context.Context, h *Handler, cfg *config.Config) *Router {
	rout := mux.NewRouter()

	logger := ctx.Value(Key).(*slog.Logger)

	rout.Use(LoggingMiddleware(logger))

	rout.HandleFunc("/api/v1/calculate", h.CalcHandler)

	return &Router{Router: rout, config: cfg}
}

func (r *Router) Run() {

	log.Fatal(http.ListenAndServe(":"+r.config.Path, r.Router))

}
