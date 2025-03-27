package storage

//go:generate mockgen -destination=../mocks/storage_mock.go -package=mocks . KvRepository
type KvRepository interface {
	AddValue(key string, value any) error
	GetValue(key string) ([]any, error)
	UpdateValue(key string, value any) error
	DeleteValue(key string) error
}
