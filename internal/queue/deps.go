package queue

import "github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"

type Service struct {
	Queue *ExpressionQueue
	Task  *TaskQueue
}

func NewService(cfg config.EnvConfig) *Service {
	return &Service{
		Queue: NewExpressionQueue(),
		Task: NewTaskQueue(Timings{
			TimeAdditionMs:       int32(cfg.TIME_ADDITION_MS),
			TimeSubtractionMs:    int32(cfg.TIME_SUBTRACTION_MS),
			TimeMultiplicationMs: int32(cfg.TIME_MULTIPLICATIONS_MS),
			TimeDivisionMs:       int32(cfg.TIME_DIVISIONS_MS),
		}),
	}
}
