package storage

const (
	JsonDataSpace string = "json_data"
	PrimaryIndex  string = "primary"

	ErrKeyNotFound string = "Key not found"
)

type TarantoolRepo interface {
	AddValue(key string, value interface{}) error
	GetValue(key string) ([]interface{}, error)
	UpdateValue(key string, value interface{}) error
	DeleteValue(key string) error
}
