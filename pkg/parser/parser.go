package parser

import (
	"github.com/ziollek/etherscription/pkg/model"
	"github.com/ziollek/etherscription/pkg/storage"
)

const (
	lastBlockKey = "last_block"
)

type SubscriptionService struct {
	txStorage    storage.ListSaver[model.Transaction]
	subStorage   storage.KVSaver[string]
	stateStorage storage.KVSaver[int]
}

func NewSubscriptionService(txStorage storage.ListSaver[model.Transaction], subStorage storage.KVSaver[string], stateStorage storage.KVSaver[int]) Parser {
	return &SubscriptionService{
		txStorage:    txStorage,
		subStorage:   subStorage,
		stateStorage: stateStorage,
	}
}

func (service *SubscriptionService) GetCurrentBlock() int {
	if block, found := service.stateStorage.Get(lastBlockKey); found {
		return block
	}
	return 0
}

func (service *SubscriptionService) Subscribe(address string) bool {
	if _, found := service.subStorage.Get(address); !found {
		service.subStorage.Set(address, address)
		return true
	}
	return false
}

func (service *SubscriptionService) IsSubscribed(address string) bool {
	_, found := service.subStorage.Get(address)
	return found
}

func (service *SubscriptionService) GetTransactions(address string) []model.Transaction {
	return service.txStorage.FetchAndFlush(address)
}
