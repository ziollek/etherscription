package model

import (
	"math/big"
	"strings"
)

// Transaction represents a simplified transaction in the Ethereum network. It is used in parser package.
type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int64  `json:"value"`
}

// RawTransaction represents a raw transaction in the Ethereum network. It is used in ethereum package.
type RawTransaction struct {
	From             string `json:"from"`
	To               string `json:"to"`
	Input            string `json:"input"`
	Value            string `json:"value"`
	TransactionIndex string `json:"transactionIndex"`
}

func (t *RawTransaction) ToTransaction() Transaction {
	return Transaction{
		From:  t.From,
		To:    t.To,
		Value: ConvertHexToInt(t.Value),
	}
}

func ConvertHexToInt(hex string) int64 {
	if !strings.HasPrefix(hex, "0x") {
		return 0
	}
	value := new(big.Int)
	value.SetString(hex, 0)
	return value.Int64()
}
