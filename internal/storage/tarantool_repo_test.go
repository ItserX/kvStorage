package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/tarantool/go-tarantool/v2"

	"kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

type Case struct {
	key           string
	value         any
	operation     func(string, any) error
	method        string
	expectedError error
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

	err = log.SetupLogger()
	if err != nil {
		t.Errorf("failed to initialize logger: %v", err)
		return
	}
	repo := storage.NewTarantoolRepository(conn)
	cases := []Case{
		{
			key:           "test",
			value:         map[string]any{"v1": 1, "v2": true, "v3": "word"},
			method:        "Add",
			expectedError: nil,
			operation:     repo.AddValue},

		{
			key:           "test",
			method:        "Get",
			expectedError: nil,
			operation: func(key string, value any) error {
				_, err := repo.GetValue(key)
				return err
			},
		},

		{
			key:           "test",
			value:         []any{1, 2, true},
			method:        "Update",
			expectedError: nil,
			operation:     repo.UpdateValue,
		},

		{
			key:           "test",
			method:        "Get",
			expectedError: nil,
			operation: func(key string, value any) error {
				_, err := repo.GetValue(key)
				return err
			},
		},

		{
			key:           "test",
			method:        "Delete",
			expectedError: nil,
			operation: func(key string, value any) error {
				err := repo.DeleteValue(key)
				return err
			},
		},

		{
			key:           "test",
			method:        "Get",
			expectedError: storage.ErrKeyNotFound,
			operation: func(key string, value any) error {
				_, err := repo.GetValue(key)
				return err
			},
		},
	}

	for _, q := range cases {
		t.Run(q.key+" "+q.method, func(t *testing.T) {
			err := q.operation(q.key, q.value)

			if q.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if err.Error() != q.expectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", q.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
