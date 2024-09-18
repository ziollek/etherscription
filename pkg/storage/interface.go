package storage

import "time"

type KVSaver[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
}

type ListSaver[T any] interface {
	FetchAndFlush(key string) []T
	Append(key string, value T, ttl time.Duration)
}
