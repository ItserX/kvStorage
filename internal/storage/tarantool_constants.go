package storage

import "errors"

const (
	JsonDataSpace string = "json_data"
	PrimaryIndex  string = "primary"
)

var ErrKeyNotFound = errors.New("key not found")
