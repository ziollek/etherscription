package memory

import (
	"context"
	"time"

	"github.com/ziollek/etherscription/pkg/logging"
)

type Cleaner[T any] struct {
	storage  *ListStorage[T]
	interval time.Duration
}

func NewCleaner[T any](storage *ListStorage[T], interval time.Duration) *Cleaner[T] {
	return &Cleaner[T]{
		storage:  storage,
		interval: interval,
	}
}

func (cleaner *Cleaner[_]) Start(ctx context.Context) {
	ticker := time.NewTicker(cleaner.interval)
	for {
		select {
		case <-ctx.Done():
			logging.Logger().Warn().
				Str("module", "memory").Msg("Context done, stopping cleaner")
		case <-ticker.C:
			totalBefore, totalAfter := 0, 0
			start := time.Now()
			for _, k := range cleaner.storage.GetKeys() {
				before, after := cleaner.storage.CleanOutdated(k)
				totalBefore += before
				totalAfter += after
			}
			logging.Logger().Info().Str("module", "memory").Dur("duration", time.Since(start)).Msgf(
				"Cleaned %d entries, left %d", totalBefore-totalAfter, totalAfter,
			)
		}
	}
}
