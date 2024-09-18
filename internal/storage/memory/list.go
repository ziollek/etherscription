package memory

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Value      T
	Expiration time.Time
}

type Entries[T any] []Entry[T]

func (entries Entries[T]) Expire(now time.Time) Entries[T] {
	newEntries := Entries[T]{}
	for _, entry := range entries {
		if entry.Expiration.After(now) {
			newEntries = append(newEntries, entry)
		}
	}
	return newEntries
}

func (entries Entries[T]) Values() []T {
	values := make([]T, len(entries))
	for i, entry := range entries {
		values[i] = entry.Value
	}
	return values
}

type ListStorage[T any] struct {
	entries map[string]Entries[T]
	sync.RWMutex
}

func NewListStorage[T any]() *ListStorage[T] {
	return &ListStorage[T]{
		entries: make(map[string]Entries[T]),
	}
}

func (storage *ListStorage[T]) FetchAndFlush(key string) []T {
	// this operation must be atomic to not lose any (not expired) data
	storage.Lock()
	defer storage.Unlock()
	if list, found := storage.entries[key]; found {
		delete(storage.entries, key)
		return list.Expire(time.Now()).Values()
	}
	return []T{}
}

func (storage *ListStorage[T]) Append(key string, value T, ttl time.Duration) {
	storage.Lock()
	defer storage.Unlock()
	if _, found := storage.entries[key]; !found {
		storage.entries[key] = Entries[T]{}
	}
	storage.entries[key] = append(storage.entries[key], Entry[T]{value, time.Now().Add(ttl)})
}

func (storage *ListStorage[_]) GetKeys() []string {
	storage.RLock()
	defer storage.RUnlock()
	keys := make([]string, 0, len(storage.entries))
	for key := range storage.entries {
		keys = append(keys, key)
	}
	return keys
}

func (storage *ListStorage[_]) CleanOutdated(key string) (int, int) {
	storage.Lock()
	defer storage.Unlock()
	if entries, found := storage.entries[key]; found {
		after := entries.Expire(time.Now())
		if len(after) > 0 {
			storage.entries[key] = after
		} else {
			// we should not store empty lists
			// because the size of the map will grow indefinitely
			delete(storage.entries, key)
		}
		return len(entries), len(after)
	}
	return 0, 0
}
