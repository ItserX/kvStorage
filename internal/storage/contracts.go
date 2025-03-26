package storage

type TarantoolRepo interface {
	AddValue(key string, value interface{}) error
	GetValue(key string) ([]interface{}, error)
	UpdateValue(key string, value interface{}) error
	DeleteValue(key string) error
}
