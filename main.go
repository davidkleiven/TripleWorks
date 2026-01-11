package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"com.github/davidkleiven/tripleworks/api"
	"com.github/davidkleiven/tripleworks/pkg"
)

func main() {
	logHandler := pkg.CtxHandler{
		Handler: slog.NewJSONHandler(os.Stdout, nil),
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	configName := os.Getenv("TRIPLE_WORKS_CONFIG")
	config := pkg.GetConfig(configName)

	mux := http.NewServeMux()
	api.Setup(mux, config)

	slog.Info("Starting server", "port", config.Port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: api.LogRequest(mux)}

	go func() {
		server.ListenAndServe()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	fmt.Println("Server stopped")
}
