package parser

import "github.com/ziollek/etherscription/pkg/model"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []model.Transaction
	IsSubscribed(address string) bool
}

type Consumer[T any] interface {
	Consume(T)
}
