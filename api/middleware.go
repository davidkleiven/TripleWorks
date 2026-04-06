package api

import (
	"context"
	"log/slog"
	"net/http"

	"com.github/davidkleiven/tripleworks/pkg"
)

func LogRequest(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), pkg.HostKey, r.RemoteAddr)
		ctx = context.WithValue(ctx, pkg.MethodKey, r.Method)
		slog.InfoContext(ctx, "Received request")
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func NoopMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
