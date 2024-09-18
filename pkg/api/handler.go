package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ziollek/etherscription/pkg/parser"
)

type Handler struct {
	parser parser.Parser
}

func NewHandler(parser parser.Parser) *Handler {
	return &Handler{parser: parser}
}

func (h *Handler) GetCurrentBlock(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	Response(w, http.StatusOK, GetCurrentBlocResponse{BlockID: h.parser.GetCurrentBlock()})
}

func (h *Handler) GetTransactions(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	if !h.parser.IsSubscribed(params.ByName("address")) {
		ErrorResponse(http.StatusNotFound, "There is no subscription for address", w)
		return
	}
	Response(w, http.StatusOK, GetTransactionsResponse{Transactions: h.parser.GetTransactions(params.ByName("address"))})
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var entry SubscriptionsRequests
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		ErrorResponse(http.StatusBadRequest, "Invalid request body", w)
		return
	}
	if h.parser.Subscribe(entry.Address) {
		Response(w, http.StatusCreated, SubscriptionsResponse{Status: true})
		return
	}
	Response(w, http.StatusOK, SubscriptionsResponse{Status: false})
}
