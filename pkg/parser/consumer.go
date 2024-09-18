package parser

import (
	"github.com/ziollek/etherscription/pkg/config"
	"github.com/ziollek/etherscription/pkg/logging"
	"github.com/ziollek/etherscription/pkg/model"
	"github.com/ziollek/etherscription/pkg/storage"
)

type TransactionConsumerService struct {
	txStorage  storage.ListSaver[model.Transaction]
	subStorage storage.KVSaver[string]
	cfg        *config.StorageConfig
}

func NewConsumerService(txStorage storage.ListSaver[model.Transaction], subStorage storage.KVSaver[string], cfg *config.StorageConfig) Consumer[model.Transaction] {
	return &TransactionConsumerService{
		txStorage:  txStorage,
		subStorage: subStorage,
		cfg:        cfg,
	}
}

func (s *TransactionConsumerService) Consume(transaction model.Transaction) {
	if _, found := s.subStorage.Get(transaction.To); found || s.cfg.StoreAllTransactions {
		logging.Logger().Debug().Str("module", "parser").Str("subscriber", transaction.To).Msgf("Appending transaction %+v to storage", transaction)
		s.txStorage.Append(transaction.To, transaction, s.cfg.Retention)
	}
	if _, found := s.subStorage.Get(transaction.From); found || s.cfg.StoreAllTransactions {
		logging.Logger().Debug().Str("module", "parser").Str("subscriber", transaction.From).Msgf("Appending transaction %+v to storage", transaction)
		s.txStorage.Append(transaction.From, transaction, s.cfg.Retention)
	}
	logging.Logger().Info().Str("module", "parser").Str("from", transaction.From).Str("to", transaction.To).Int64("value", transaction.Value).Msgf("consumed")
}

type StateConsumerService struct {
	stateStorage storage.KVSaver[int]
}

func NewStateConsumerService(stateStorage storage.KVSaver[int]) Consumer[int] {
	return &StateConsumerService{
		stateStorage: stateStorage,
	}
}

func (s *StateConsumerService) Consume(lastBlock int) {
	logging.Logger().Info().Str("module", "parser").Int(lastBlockKey, lastBlock).Msgf("Updating last block")
	s.stateStorage.Set(lastBlockKey, lastBlock)
}
