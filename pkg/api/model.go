package api

import "github.com/ziollek/etherscription/pkg/model"

type JSONErrorResponse struct {
	Error *GenericError `json:"error"`
}

type GenericError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type SubscriptionsRequests struct {
	Address string `json:"address"`
}

type SubscriptionsResponse struct {
	Status bool `json:"status"`
}

type GetCurrentBlocResponse struct {
	BlockID int `json:"block_id"`
}

type GetTransactionsResponse struct {
	Transactions []model.Transaction `json:"transactions"`
}
