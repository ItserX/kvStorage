package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	tarantool "github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"kvManager/internal/handlers"
	"kvManager/internal/storage"
)

func connectToTarantool(addr string, user string, logger *zap.SugaredLogger) (*tarantool.Connection, error) {
	logger.Infow("Connecting to Tarantool", "address", addr, "user", user)

	dialer := tarantool.NetDialer{
		Address: addr,
		User:    user,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := tarantool.Opts{
		Timeout: 5 * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		logger.Errorw("Failed to connect to Tarantool", "error", err, "address", addr)
		return nil, err
	}
	defer conn.Close()

	logger.Info("Successfully connected to Tarantool")
	return conn, nil
}

func setupRouter(conn *tarantool.Connection, logger *zap.SugaredLogger) *mux.Router {
	logger.Info("Setting up router and initializing storage")

	st := storage.NewTarantoolRepository(conn, logger)
	h := handlers.Handler{Repo: st, Logger: logger}

	r := mux.NewRouter()
	r.HandleFunc("/kv", h.Add).Methods("POST")
	r.HandleFunc("/kv/{id}", h.Get).Methods("GET")
	r.HandleFunc("/kv/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/kv/{id}", h.Delete).Methods("DELETE")

	logger.Info("Router setup completed")
	return r
}

func setupLogger() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableStacktrace = true

	logger, err := config.Build()

	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return sugar, nil
}

func main() {
	logger, err := setupLogger()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		return
	}

	logger.Info("Starting app")
	conn, err := connectToTarantool("tarantool:3301", "guest", logger)
	if err != nil {
		return
	}

	r := setupRouter(conn, logger)

	logger.Infow("Starting HTTP server", "address", ":8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Errorw("HTTP server error", err)
		return
	}
}
