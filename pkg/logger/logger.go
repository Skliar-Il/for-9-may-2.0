package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type KeyLoggerType string

const (
	RequestId KeyLoggerType = "request_id"
	lKey      KeyLoggerType = "logger"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx *gin.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	ctx.Set(string(lKey), &Logger{logger})
	return nil
}

func GetLoggerFromCtx(ctx *gin.Context) *Logger {
	logger, exist := ctx.Get(string(lKey))
	if !exist {
		return nil
	}
	return logger.(*Logger)
}

func (l *Logger) Info(ctx *gin.Context, msg string, fields ...zap.Field) {
	if ctx.GetString(string(RequestId)) != "" {
		fields = append(fields, zap.String(string(RequestId), ctx.GetString(string(RequestId))))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(ctx *gin.Context, msg string, fields ...zap.Field) {
	if ctx.GetString(string(RequestId)) != "" {
		fields = append(fields, zap.String(string(RequestId), ctx.GetString(string(RequestId))))
	}
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx *gin.Context, msg string, fields ...zap.Field) {
	if ctx.GetString(string(RequestId)) != "" {
		fields = append(fields, zap.String(string(RequestId), ctx.GetString(string(RequestId))))
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
