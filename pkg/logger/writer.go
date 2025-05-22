package logger

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
)

type otelWriter struct {
	span trace.Span
}

func (w *otelWriter) Write(p []byte) (n int, err error) {
	w.span.AddEvent(string(p))
	fmt.Println("huy")
	return len(p), nil
}
