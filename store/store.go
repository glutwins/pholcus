package store

import (
	"github.com/glutwins/pholcus/config"
)

type Storage interface {
	FetchKVData(string) (map[string]interface{}, error)
	InsertStringMap(string, map[string]interface{}) error
	InsertKVData(string, map[string]interface{}) (int, error)
	ClearKVData(string) error
	Close() error
}

func NewStorage(db *config.PholcusDbConfig) Storage {
	return nil
}
