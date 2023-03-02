package history

import (
	"KeySell/pkg/auth"
	"KeySell/pkg/tooling"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type HistoryService interface {
	GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error)
	GetOneTransaction(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error)
	EditTransaction(ctx context.Context, UserID, TransactID int, Key string) error
}

type HistoryHandler struct {
	ware auth.MiddleWare
	hs   HistoryService
}

func NewHistoryHandler(ware auth.MiddleWare, hs HistoryService) *HistoryHandler {
	return &HistoryHandler{ware: ware, hs: hs}
}

func (h *HistoryHandler) Register(r *httprouter.Router) {
	r.GET("/api/seller/history", h.ware.IsAuth(h.GetHistory))
	r.GET("/api/seller/history/:tran_id", h.ware.IsAuth(h.GetFullTransaction))
	r.PATCH("/api/seller/history/:tran_id", h.ware.IsAuth(h.EditTransaction))
}

func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	transactions, err := h.hs.GetAllTransactions(r.Context(), UserID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}

	transactionsMarshalled, err := json.Marshal(transactions)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write(transactionsMarshalled)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *HistoryHandler) GetFullTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	tranID := tooling.GetTranID(params)

	transaction, err := h.hs.GetOneTransaction(r.Context(), UserID, tranID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	marshal, err := json.Marshal(transaction)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

type EditTransationInput struct {
	ContentKey string `json:"content_key,omitempty"`
}

func (h *HistoryHandler) EditTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	tranID := tooling.GetTranID(params)

	var input EditTransationInput

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		logrus.Println(err)
		return
	}

	err = h.hs.EditTransaction(r.Context(), UserID, tranID, input.ContentKey)
	if err != nil {
		logrus.Println(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
		return
	}

}
