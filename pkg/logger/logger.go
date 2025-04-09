package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type KeyLoggerType string

const (
	RequestId KeyLoggerType = "request_id"
	lKey      KeyLoggerType = "logger"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	if ctx.Value(RequestId) == nil {
		ctx = context.WithValue(ctx, RequestId, uuid.New().String())
	}

	ctx = context.WithValue(ctx, lKey, &Logger{logger})
	return nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	logger, exist := ctx.Value(lKey).(*Logger)
	if !exist {
		return nil
	}
	return logger
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if id, ok := ctx.Value(RequestId).(string); ok {
		fields = append(fields, zap.String(string(RequestId), id))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if id, ok := ctx.Value(RequestId).(string); ok {
		fields = append(fields, zap.String(string(RequestId), id))
	}
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if id, ok := ctx.Value(RequestId).(string); ok {
		fields = append(fields, zap.String(string(RequestId), id))
	}
	l.l.Fatal(msg, fields...)
}

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		guid := uuid.New().String()
		ctx.Set(string(RequestId), guid)

		if GetLoggerFromCtx(ctx) == nil {
			_ = New(ctx)
		}
		GetLoggerFromCtx(ctx).Info(ctx,
			"Request http",
			zap.String("method", ctx.Request.Method),
		)

		timeStart := time.Now()
		ctx.Next()
		duration := time.Since(timeStart)

		GetLoggerFromCtx(ctx).Info(ctx,
			"Response http",
			zap.Duration("duration", duration),
		)
	}
}
