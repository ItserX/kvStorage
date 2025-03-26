package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"kvManager/pkg/handlers"
	"kvManager/pkg/storage"
)

type Case struct {
	Path           string
	Body           string
	Method         string
	ExpectedStatus int
}

func TestAPIHandlers(t *testing.T) {

	dialer := tarantool.NetDialer{
		Address: ":3301",
		User:    "guest",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := tarantool.Opts{
		Timeout: 5 * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		panic(err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableStacktrace = true

	logger, _ := config.Build()
	defer logger.Sync()

	sugar := logger.Sugar()

	repo := storage.NewTarantoolRepository(conn, sugar)
	handler := handlers.DbHandler{Repo: repo, Logger: sugar}

	router := mux.NewRouter()
	router.HandleFunc("/kv", handler.Add).Methods("POST")
	router.HandleFunc("/kv/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/kv/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/kv/{id}", handler.Delete).Methods("DELETE")

	ts := httptest.NewServer(router)
	defer ts.Close()

	Cases := []Case{
		{
			Path:           "/kv",
			Body:           `{"key":"test1", "value":{"k1":123, "k2":true}}`,
			Method:         "POST",
			ExpectedStatus: http.StatusCreated,
		},
		{
			Path:           "/kv",
			Body:           `{"key": "test2", "value": [1, true, "word"]}`,
			Method:         "POST",
			ExpectedStatus: http.StatusCreated,
		},
		{
			Path: "/kv",
			Body: `{
				  "key": "test3",
				  "value": {
					"user": {
					  "id": 1,
					  "profile": {
						"name": "name",
						"settings": {
						  "theme": "dark",
						  "notifications": {
							"email": true,
							"push": false
						  }
						}
					  }
					},
					"system": {
					  "version": "1.2.3",
					  "dependencies": {
						"db": "tarantool"
					  }
					}
				  }
				}`,
			Method:         "POST",
			ExpectedStatus: http.StatusCreated,
		},
		{
			Path:           "/kv",
			Body:           `{"key": "test2", "value": [1, true, "word"]}`,
			Method:         "POST",
			ExpectedStatus: http.StatusConflict,
		},
		{
			Path:           "/kv",
			Body:           `{"key": "test4", {inccorect body}`,
			Method:         "POST",
			ExpectedStatus: http.StatusBadRequest,
		},

		{
			Path:           "/kv/test1",
			Method:         "GET",
			ExpectedStatus: http.StatusOK,
		},
		{
			Path:           "/kv/test1",
			Body:           `{"value":"new_value"}`,
			Method:         "PUT",
			ExpectedStatus: http.StatusOK,
		},
		{
			Path:           "/kv/test1",
			Body:           `incorrect_body`,
			Method:         "PUT",
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Path:           "/kv/invalid",
			Body:           `{"value":"new_value"}`,
			Method:         "PUT",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Path:           "/kv/test1",
			Method:         "DELETE",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Path:           "/kv/test1",
			Method:         "DELETE",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Path:           "/kv/test2",
			Method:         "DELETE",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Path:           "/kv/test3",
			Method:         "DELETE",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Path:           "/kv/non",
			Method:         "GET",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Path:           "/kv",
			Body:           `incorrect_body`,
			Method:         "POST",
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for idx, q := range Cases {
		t.Run(q.Path, func(t *testing.T) {
			url := ts.URL + q.Path

			caseName := fmt.Sprintf("case %d: %s %s", idx, q.Method, q.Path)

			var req *http.Request
			var err error

			if q.Body != "" {
				req, err = http.NewRequest(q.Method, url, bytes.NewBufferString(q.Body))
				if err != nil {
					t.Fatalf("[%s] error: %v", caseName, err)
				}
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(q.Method, url, nil)
				if err != nil {
					t.Fatalf("[%s] error: %v", caseName, err)
				}
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("[%s] error: %v", caseName, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != q.ExpectedStatus {
				t.Errorf("Expected status %d, got %d", q.ExpectedStatus, resp.StatusCode)
			}

		})
	}
}
