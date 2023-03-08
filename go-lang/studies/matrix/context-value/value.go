package context_value

import (
	"context"
	"net/http"
)

type traceId struct{}

func ParseTraceId(ctx context.Context) string {
	if value, ok := ctx.Value(traceId{}).(string); ok {
		return value
	}
	return ""
}

func WithTraceId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, traceId{}, value)
}

type verbose struct{}

func ParseVerbose(ctx context.Context) bool {
	if value, ok := ctx.Value(verbose{}).(bool); ok {
		return value
	}
	return false
}

func WithVerbose(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, verbose{}, value)
}

type request struct{}

func ParseRequest(ctx context.Context) *http.Request {
	if value, ok := ctx.Value(request{}).(*http.Request); ok {
		return value
	}
	return nil
}

func WithRequest(ctx context.Context, value *http.Request) context.Context {
	return context.WithValue(ctx, request{}, value)
}

type responseWriter struct{}

func ParseResponseWriter(ctx context.Context) http.ResponseWriter {
	if value, ok := ctx.Value(responseWriter{}).(http.ResponseWriter); ok {
		return value
	}
	return nil
}

func WithResponseWriter(ctx context.Context, value http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseWriter{}, value)
}
