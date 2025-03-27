package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	tarantool "github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"

	"kvManager/internal/handlers"
	log "kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func connectToTarantool(addr string, user string) (*tarantool.Connection, error) {
	log.Logger.Infow("Connecting to Tarantool", "address", addr, "user", user)

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
		log.Logger.Errorw("Failed to connect to Tarantool", "error", err, "address", addr)
		return nil, err
	}

	log.Logger.Info("Successfully connected to Tarantool")
	return conn, nil
}

func setupRouter(conn *tarantool.Connection, logger *zap.SugaredLogger) *mux.Router {
	logger.Info("Setting up router and initializing storage")

	st := storage.NewTarantoolRepository(conn)
	h := handlers.Handler{Repo: st}

	r := mux.NewRouter()
	r.HandleFunc("/kv", h.Add).Methods("POST")
	r.HandleFunc("/kv/{id}", h.Get).Methods("GET")
	r.HandleFunc("/kv/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/kv/{id}", h.Delete).Methods("DELETE")

	logger.Info("Router setup completed")
	return r
}

func main() {
	err := log.SetupLogger()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		return
	}

	err = loadEnv()
	if err != nil {
		log.Logger.Errorw("Env load failing")
		return
	}

	appPort := os.Getenv("APP_PORT")
	tarantoolAddr := os.Getenv("TARANTOOL_ADDRESS")
	tarantoolUser := os.Getenv("TARANTOOL_USER")

	log.Logger.Info("Starting app")
	conn, err := connectToTarantool(tarantoolAddr, tarantoolUser)
	if err != nil {
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Logger.Errorw("Connection to tarantool is not closed", err)
		}
	}()

	r := setupRouter(conn, log.Logger)

	log.Logger.Infow("Starting HTTP server", "address", appPort)
	err = http.ListenAndServe(appPort, r)
	if err != nil {
		log.Logger.Errorw("HTTP server error", err)
		return
	}
}
