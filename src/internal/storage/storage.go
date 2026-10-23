package storage

import (
	"sync"
)

type StorageValue struct {
	DataType string
	Data     string
}

type StorageMap struct {
	mu   sync.RWMutex
	data map[string]StorageValue
}

func (storage *StorageMap) Put(key string, value StorageValue) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	if storage.data == nil {
		storage.data = make(map[string]StorageValue)
	}
	storage.data[key] = value
}
func (storage *StorageMap) Get(key string) (StorageValue, bool) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	value, ok := storage.data[key]
	return value, ok
}
func (storage *StorageMap) Clean(key string) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	delete(storage.data, key)
}
