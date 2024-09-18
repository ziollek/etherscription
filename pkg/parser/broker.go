package parser

import (
	"context"

	"github.com/ziollek/etherscription/pkg/logging"
	"github.com/ziollek/etherscription/pkg/model"
)

type Broker struct {
	transactions  <-chan model.Transaction
	blocks        <-chan int
	txConsumer    Consumer[model.Transaction]
	stateConsumer Consumer[int]
}

func NewBroker(transactions <-chan model.Transaction, blocks <-chan int, txConsumer Consumer[model.Transaction], stateConsumer Consumer[int]) *Broker {
	return &Broker{
		transactions:  transactions,
		blocks:        blocks,
		txConsumer:    txConsumer,
		stateConsumer: stateConsumer,
	}
}

func (broker *Broker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logging.Logger().Warn().
				Str("module", "parser").Msg("Context done, stopping fetcher")
			return
		case transaction := <-broker.transactions:
			// it can be done in parallel for slower storages
			broker.txConsumer.Consume(transaction)
		case block := <-broker.blocks:
			broker.stateConsumer.Consume(block)
		}
	}
}
