package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"kvManager/pkg/storage"
)

type Case struct {
	Key           string
	Value         interface{}
	Method        string
	ExpectedError error
}

func TestTarantoolRepo(t *testing.T) {
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
		t.Fatalf("Failed to connect: %v", err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, _ := config.Build()
	defer logger.Sync()

	sugar := logger.Sugar()

	repo := storage.NewTarantoolRepository(conn, sugar)
	cases := []Case{
		{Key: "test", Value: map[string]interface{}{"v1": 1, "v2": true, "v3": "word"}, Method: "Add", ExpectedError: nil},
		{Key: "test", Method: "Get", ExpectedError: nil},
		{Key: "test", Value: []interface{}{1, 2, true}, Method: "Update", ExpectedError: nil},
		{Key: "test", Method: "Get", ExpectedError: nil},
		{Key: "test", Method: "Delete", ExpectedError: nil},
		{Key: "test", Method: "Get", ExpectedError: fmt.Errorf(storage.ErrKeyNotFound)},
	}

	for _, q := range cases {
		t.Run(q.Key+" "+q.Method, func(t *testing.T) {
			var err error

			switch q.Method {
			case "Get":
				_, err = repo.GetValue(q.Key)
			case "Add":
				err = repo.AddValue(q.Key, q.Value)
			case "Update":
				err = repo.UpdateValue(q.Key, q.Value)
			case "Delete":
				err = repo.DeleteValue(q.Key)
			default:
				t.Fatalf("Unknown method: %s", q.Method)
			}

			if q.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if err.Error() != q.ExpectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", q.ExpectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
