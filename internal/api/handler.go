package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/queue"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg/logger"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Orchestrator struct {
	mu          sync.Mutex
	expressions map[string]*domain.Expression

	cfg *config.EnvConfig

	s *queue.Service
}

func NewOrchestrator(cfg config.EnvConfig, s *queue.Service) *Orchestrator {
	return &Orchestrator{
		expressions: make(map[string]*domain.Expression),
		s:           s,
		mu:          sync.Mutex{},
		cfg:         &cfg,
	}
}

func (o *Orchestrator) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req domain.JsonReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	expr := &domain.Expression{ID: utils.IdGen(), Status: "pending"}

	o.s.Queue.AddExpression(expr)

	response := map[string]string{"id": expr.ID}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Println("lets calculete")

}

func (o *Orchestrator) ExpressionsHandel(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/expression/" {

		expr := o.s.Queue.GetExpressions()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expr)

		return

	}

	id := r.URL.Path[len("/api/v1/expression/"):]

	expr, ok := o.s.Queue.GetExpressionById(id)

	if !ok {

		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(expr)

}

func (o *Orchestrator) InternalTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		task, ok := o.s.Task.GetTask()
		if !ok {
			http.Error(w, "Not task", http.StatusNotFound)
		}

		bytew, err := json.Marshal(task)
		if err != nil {
			return
		}
		fmt.Fprint(w, bytew)
		return
	}
	var input domain.Expression

	json.NewDecoder(r.Body).Decode(&input)

	err := o.s.Task.RollbackResult(input.ID, *input.Result)

	if err != nil {
		http.Error(w, "Error server", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

// Мидлвaря няшка
func LoggingMiddleware(ctx context.Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "Finished:",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", http.StatusOK),
				zap.Duration("duration", time.Second),
			)

			next.ServeHTTP(w, r)

		})
	}
}
