package etherum

import "github.com/ziollek/etherscription/pkg/model"

type logEntry struct {
	Removed         bool   `json:"removed"`
	BlockNumber     string `json:"blockNumber"`
	TransactionHash string `json:"transactionHash"`
}

type logEntries []logEntry

func (entries logEntries) GetUniqueTransactionHashes() []string {
	uniqueTransactionHashes := make(map[string]struct{})
	for _, entry := range entries {
		if !entry.Removed {
			uniqueTransactionHashes[entry.TransactionHash] = struct{}{}
		}
	}
	var result []string
	for hash := range uniqueTransactionHashes {
		result = append(result, hash)
	}
	return result
}

func (entries logEntries) GetLastBlock() int {
	block := "0x0"
	for _, entry := range entries {
		if entry.BlockNumber != block {
			block = entry.BlockNumber
		}
	}
	return int(model.ConvertHexToInt(block))
}
