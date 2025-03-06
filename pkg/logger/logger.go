package logger

import (
	"context"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/domain"
	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, domain.Logger, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(domain.Logger).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(domain.RequestID) != nil {
		fields = append(fields, zap.String(string(domain.RequestID), ctx.Value(domain.RequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(domain.RequestID) != nil {
		fields = append(fields, zap.String(string(domain.RequestID), ctx.Value(domain.RequestID).(string)))
	}
	l.l.Fatal(msg, fields...)
}
