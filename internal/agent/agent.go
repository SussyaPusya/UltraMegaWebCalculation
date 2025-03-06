package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg/logger"
	"go.uber.org/zap"
)

type Agent struct {
	cfg *config.Config
	mu  sync.Mutex
}

func NewAgent(cfg *config.Config) *Agent {
	return &Agent{cfg: cfg, mu: sync.Mutex{}}
}

func (a *Agent) Run(ctx context.Context) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "agent started")

	computingPower := a.cfg.EnvConfig.COMPUTING_POWER
	orchestratorPort := a.cfg.Port
	if computingPower <= 0 {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "invalid computing power")
		os.Exit(1)
	}
	for i := 0; i < computingPower; i++ {
		go worker(orchestratorPort, ctx)
		time.Sleep(time.Millisecond * 10)
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("created worker %d", i+1))
	}
}

func worker(port string, ctx context.Context) {
	for {
		task, err := fetchTask(port)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		result, err := solveTask(task, ctx)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to solve task:", zap.Error(err))
			continue
		}

		err = sendResult(port, task.ID, result)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to send result:", zap.Error(err))
			continue
		}

		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("solved task %s: result=%f", task.ID, result))
	}
}

func fetchTask(port string) (*domain.Task, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/internal/task", port))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch task: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var task domain.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func solveTask(task *domain.Task, ctx context.Context) (float64, error) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Solving task %2.f%s%2.f", task.Arg1, task.Operation, task.Arg2))
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", task.Operation)
	}
}

func sendResult(port, taskId string, result float64) error {
	data := map[string]interface{}{
		"id":     taskId,
		"result": result,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/internal/task", port), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send result: %s", string(body))
	}

	return nil
}
