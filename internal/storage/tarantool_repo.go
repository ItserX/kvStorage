package storage

import (
	"fmt"
	
	"github.com/tarantool/go-tarantool/v2"

	"kvManager/internal/pkg/log"
)

type TarantoolRepository struct {
	conn *tarantool.Connection
}

func NewTarantoolRepository(conn *tarantool.Connection) *TarantoolRepository {
	return &TarantoolRepository{conn: conn}
}

func (repo *TarantoolRepository) execRequest(req tarantool.Request) ([]any, error) {
	future := repo.conn.Do(req)
	data, err := future.Get()
	if err != nil {
		log.Logger.Warnw("Duplicate key error",
			"error", err.Error())
		return nil, err
	}

	log.Logger.Debugw("Tarantool response data",
		"data_length", len(data),
		"data", fmt.Sprintf("%s", data))

	if len(data) == 0 {
		log.Logger.Debugw("Empty response from Tarantool", "operation", "execRequest")
		return nil, ErrKeyNotFound
	}
	return data, nil
}

func (repo *TarantoolRepository) AddValue(key string, value any) error {
	log.Logger.Debugw("Adding value to Tarantool",
		"key", key)
	req := tarantool.NewInsertRequest(JsonDataSpace).Tuple([]any{key, value})
	_, err := repo.execRequest(req)
	return err
}

func (repo *TarantoolRepository) GetValue(key string) ([]any, error) {
	log.Logger.Debugw("Get value from Tarantool",
		"key", key)
	req := tarantool.NewSelectRequest(JsonDataSpace).Index(PrimaryIndex).Key([]any{key})
	return repo.execRequest(req)
}

func (repo *TarantoolRepository) UpdateValue(key string, value any) error {
	log.Logger.Debugw("Update value in Tarantool",
		"key", key)
	req := tarantool.NewUpdateRequest(JsonDataSpace).Index(PrimaryIndex).Key([]any{key}).Operations(tarantool.NewOperations().Assign(1, value))
	_, err := repo.execRequest(req)
	return err
}

func (repo *TarantoolRepository) DeleteValue(key string) error {
	log.Logger.Debugw("Delete value from Tarantool",
		"key", key)
	req := tarantool.NewDeleteRequest(JsonDataSpace).Index(PrimaryIndex).Key([]any{key})
	_, err := repo.execRequest(req)
	return err
}
