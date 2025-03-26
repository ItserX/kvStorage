package storage

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

type TarantoolRepository struct {
	conn   *tarantool.Connection
	logger *zap.SugaredLogger
}

func NewTarantoolRepository(conn *tarantool.Connection, logger *zap.SugaredLogger) *TarantoolRepository {
	return &TarantoolRepository{conn: conn, logger: logger}
}

func (repo *TarantoolRepository) execRequest(req tarantool.Request) ([]interface{}, error) {
	future := repo.conn.Do(req)
	data, err := future.Get()
	if err != nil {
		repo.logger.Warnw("Duplicate key error",
			"error", err.Error())
		return nil, err
	}

	repo.logger.Debugw("Tarantool response data",
		"data_length", len(data),
		"data", fmt.Sprintf("%s", data))

	if len(data) == 0 {
		repo.logger.Debugw("Empty response from Tarantool", "operation", "execRequest")
		return nil, fmt.Errorf("%s", ErrKeyNotFound)
	}

	return data, nil
}

func (repo *TarantoolRepository) AddValue(key string, value interface{}) error {
	repo.logger.Debugw("Adding value to Tarantool",
		"key", key)
	req := tarantool.NewInsertRequest(JsonDataSpace).Tuple([]interface{}{key, value})
	_, err := repo.execRequest(req)
	return err
}

func (repo *TarantoolRepository) GetValue(key string) ([]interface{}, error) {
	repo.logger.Debugw("Get value from Tarantool",
		"key", key)
	req := tarantool.NewSelectRequest(JsonDataSpace).Index(PrimaryIndex).Key([]interface{}{key})
	return repo.execRequest(req)
}

func (repo *TarantoolRepository) UpdateValue(key string, value interface{}) error {
	repo.logger.Debugw("Update value in Tarantool",
		"key", key)
	req := tarantool.NewUpdateRequest(JsonDataSpace).Index(PrimaryIndex).Key([]interface{}{key}).Operations(tarantool.NewOperations().Assign(1, value))
	_, err := repo.execRequest(req)
	return err
}

func (repo *TarantoolRepository) DeleteValue(key string) error {
	repo.logger.Debugw("Delete value from Tarantool",
		"key", key)
	req := tarantool.NewDeleteRequest(JsonDataSpace).Index(PrimaryIndex).Key([]interface{}{key})
	_, err := repo.execRequest(req)
	return err
}
