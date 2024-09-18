package memory

import (
	"sync"

	"github.com/ziollek/etherscription/pkg/storage"
)

type KVStorage[T any] struct {
	entries map[string]T
	sync.RWMutex
}

func NewKVStorage[T any]() storage.KVSaver[T] {
	return &KVStorage[T]{
		entries: make(map[string]T),
	}
}

func (storage *KVStorage[T]) Get(key string) (T, bool) {
	storage.RLock()
	defer storage.RUnlock()
	value, found := storage.entries[key]
	return value, found
}

func (storage *KVStorage[T]) Set(key string, value T) {
	storage.Lock()
	defer storage.Unlock()
	storage.entries[key] = value
}
