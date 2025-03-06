package api

import (
	"context"
	"log"
	"net/http"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	"github.com/gorilla/mux"
)

type key string

const Key = key("Logger")

type Router struct {
	Router  *mux.Router
	Handler Orchestrator

	config *config.Config
}

func NewRouter(ctx context.Context, h *Orchestrator, cfg *config.Config) *Router {
	rout := mux.NewRouter()

	rout.Use(LoggingMiddleware(ctx))

	rout.HandleFunc("/api/v1/calculate", h.handleCalculate)

	rout.HandleFunc("/api/v1/expressions", h.ExpressionsHandel)

	rout.HandleFunc("/internal/task", h.InternalTasks)

	return &Router{Router: rout, config: cfg}
}

func (r *Router) Run() {

	log.Fatal(http.ListenAndServe(":"+r.config.Port, r.Router))

}
