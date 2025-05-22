package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
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
	l.logWithSpan(ctx, msg, fields, l.l.Info)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSpan(ctx, msg, fields, l.l.Error)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSpan(ctx, msg, fields, l.l.Fatal)
}

func (l *Logger) logWithSpan(
	ctx context.Context,
	msg string,
	fields []zap.Field,
	logFunc func(string, ...zap.Field),
) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		logFunc(msg, fields...)
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(fields))
	for _, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			attrs = append(attrs, attribute.String(field.Key, field.String))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type,
			zapcore.Float32Type, zapcore.Float64Type,
			zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type,
			zapcore.TimeType, zapcore.DurationType:
			attrs = append(attrs, attribute.Int64(field.Key, field.Integer))
		default:
			attrs = append(attrs, attribute.String(field.Key, fmt.Sprint(field.Interface)))
		}
	}

	span.AddEvent(msg, trace.WithAttributes(attrs...))

	logFunc(msg, fields...)
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

		logger.Info(ctx.Request.Context(),
			"Request started",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
		)

		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime)

		logger.Info(ctx.Request.Context(),
			"Request completed",
			zap.Duration("duration", duration),
			zap.Int("status", ctx.Writer.Status()),
		)
	}
}
