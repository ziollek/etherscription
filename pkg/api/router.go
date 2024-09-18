package api

import (
	"github.com/julienschmidt/httprouter"
)

func ConfigureRouting(handler *Handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/current-block", handler.GetCurrentBlock)
	router.GET("/api/new-transactions/:address", handler.GetTransactions)
	router.POST("/api/subscribe", handler.Subscribe)
	return router
}
