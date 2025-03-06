package queue

import (
	"errors"
	"go/ast"

	"strconv"
	"sync"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg"
	"github.com/google/uuid"
)

type Timings struct {
	TimeAdditionMs       int32
	TimeSubtractionMs    int32
	TimeMultiplicationMs int32
	TimeDivisionMs       int32
}

type ExpressionQueue struct {
	mu          sync.RWMutex
	Expressions map[string]*domain.Expression
}

func NewExpressionQueue() *ExpressionQueue {
	return &ExpressionQueue{
		mu:          sync.RWMutex{},
		Expressions: make(map[string]*domain.Expression),
	}
}

// AddExpression with mu
func (eq *ExpressionQueue) AddExpression(expr *domain.Expression) {
	eq.mu.Lock()
	defer eq.mu.Unlock()
	eq.Expressions[expr.ID] = expr
}

// GetExpressionById with mu
func (eq *ExpressionQueue) GetExpressionById(id string) (*domain.Expression, bool) {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	expr, ok := eq.Expressions[id]
	return expr, ok
}

// GetExpressions with mu
func (eq *ExpressionQueue) GetExpressions() []domain.Expression {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	var exprs []domain.Expression
	for _, expr := range eq.Expressions {
		exprs = append(exprs, *expr)
	}
	return exprs
}

// GetExpression with mu
func (eq *ExpressionQueue) GetExpression() (*domain.Expression, bool) {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	if len(eq.Expressions) == 0 {
		return nil, false
	}
	for _, expr := range eq.Expressions {
		return expr, true
	}
	return nil, false
}

func (eq *ExpressionQueue) RemoveExpression(id string) {
	eq.mu.Lock()
	defer eq.mu.Unlock()
	delete(eq.Expressions, id)
}

type TaskQueue struct {
	mu    sync.RWMutex
	Tasks map[string]*domain.Task
	Timings
}

func NewTaskQueue(t Timings) *TaskQueue {
	return &TaskQueue{
		mu:      sync.RWMutex{},
		Tasks:   make(map[string]*domain.Task),
		Timings: t,
	}
}

// AddTask with mu
func (tq *TaskQueue) AddTask(task *domain.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.Tasks[task.ID] = task
}

// GetTask with mu
func (tq *TaskQueue) GetTask() (*domain.Task, bool) {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	if len(tq.Tasks) == 0 {
		return nil, false
	}
	for _, task := range tq.Tasks {
		return task, true
	}
	return nil, false
}

// RemoveTask with mu
func (tq *TaskQueue) RemoveTask(id string) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	delete(tq.Tasks, id)
}

// RunTask with mu
func (tq *TaskQueue) RunTask(eq *ExpressionQueue, expression *domain.Expression) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	expression.Status = "in progress"
	resChan := make(chan float64)

	astNode, err := pkg.ParseAst(expression.Expression)
	if err != nil {
		return errors.New(err.Error())
	}

	go func() {
		defer close(resChan)
		result := tq.evaluateAst(astNode, resChan)
		eq.WriteResultToExpression(expression, result)
	}()

	return nil
}

func (eq *ExpressionQueue) WriteResultToExpression(expression *domain.Expression, result float64) {
	eq.mu.Lock()
	defer eq.mu.Unlock()
	expression.Status = "Completed"
	*expression.Result = result
}

// evaluateAst with mu
func (tq *TaskQueue) evaluateAst(node ast.Node, res chan float64) float64 {
	switch n := node.(type) {
	case *ast.ParenExpr:
		return tq.evaluateAst(n.X, res)
	case *ast.BinaryExpr:
		left := tq.evaluateAst(n.X, res)
		right := tq.evaluateAst(n.Y, res)

		task := &domain.Task{
			ID:        uuid.NewString(),
			Arg1:      left,
			Arg2:      right,
			Operation: n.Op.String(),
		}
		tq.AddTask(task)

		result := <-res
		tq.RemoveTask(task.ID)
		return result
	case *ast.BasicLit:
		val, _ := strconv.ParseFloat(n.Value, 64)
		return val
	default:
		return 0.0
	}
}

// RollbackResult with mu
func (tq *TaskQueue) RollbackResult(id string, res float64) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	task, ok := tq.Tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	task.ResultChan <- res
	return nil
}
