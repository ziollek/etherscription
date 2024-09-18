package etherum

import (
	"context"
	"time"

	"github.com/ziollek/etherscription/pkg/logging"
	"github.com/ziollek/etherscription/pkg/model"
)

type Fetcher struct {
	interval   time.Duration
	client     *RPCClient
	lastBlock  int
	txChan     chan<- model.Transaction
	blocksChan chan<- int
}

func NewFetcher(interval time.Duration, client *RPCClient, txChan chan<- model.Transaction, blocksChan chan<- int) *Fetcher {
	return &Fetcher{
		lastBlock:  0,
		interval:   interval,
		client:     client,
		txChan:     txChan,
		blocksChan: blocksChan,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	// producer should close channels
	defer close(f.txChan)
	defer close(f.blocksChan)
	filterID, err := f.client.createFilter()
	if err != nil {
		logging.Logger().Err(err).Str("module", "etherum").Msg("Cannot create filter")
		return err
	}
	logging.Logger().Info().Str("module", "etherum").Msgf("Filter %s created", filterID)

	ticker := time.NewTicker(f.interval)
	for {
		select {
		case <-ctx.Done():
			logging.Logger().Warn().
				Str("module", "fetcher").Msg("Context done, stopping fetcher")
			return nil
		case <-ticker.C:
			entries, err := f.client.getChanges(filterID)
			start := time.Now()
			if err != nil {
				logging.Logger().Err(err).Str("module", "etherum").Msg("Cannot get filter changes")
			} else {
				transactions := entries.GetUniqueTransactionHashes()
				logging.Logger().Info().Str("module", "etherum").Msgf("Fetched %d uniq transactions", len(transactions))
				for _, txHash := range transactions {
					transaction, err := f.client.getTransaction(txHash)
					if err != nil {
						// there should be a retry mechanism
						logging.Logger().Err(err).Str("module", "etherum").Msgf("Cannot get transaction %s", txHash)
					} else {
						logging.Logger().Debug().
							Str("module", "etherum").
							Str("transaction", txHash).
							Msgf("New transaction fetched: %+v", transaction)
						f.txChan <- transaction.ToTransaction()
					}
				}
				if nextBlock := entries.GetLastBlock(); nextBlock > f.lastBlock {
					f.lastBlock = nextBlock
					f.blocksChan <- f.lastBlock
					logging.Logger().Debug().Str("module", "etherum").Dur("duration", time.Since(start)).Msgf("Block %d has been parsed", f.lastBlock)
				}
			}
		}
	}
}
