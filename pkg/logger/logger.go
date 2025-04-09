package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"strings"
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

func New() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Logger{l: logger}
}

func GetLoggerFromCtx(ctx *gin.Context) *Logger {
	if logger, exists := ctx.Get(string(lKey)); exists {
		return logger.(*Logger)
	}
	return nil
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
func Middleware(logger *Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/swagger/") {
			ctx.Next()
			return
		}

		guid := uuid.New().String()
		ctx.Set(string(RequestId), guid)

		ctx.Set(string(lKey), logger)

		logger.Info(ctx,
			"Request http",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
		)

		timeStart := time.Now()
		ctx.Next()
		duration := time.Since(timeStart)

		logger.Info(ctx,
			"Response http",
			zap.Duration("duration", duration),
			zap.Int("status", ctx.Writer.Status()),
		)
	}
}
