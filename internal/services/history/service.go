package history

import "context"

type HistoryStorage interface {
	GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error)
	GetOneTransaction(ctx context.Context, TransactID int) (map[string]interface{}, error)
	EditTransaction(ctx context.Context, TransactID int, Key string) error
	DeleteTransaction(ctx context.Context, TransactID int) error
}

type HistoryService struct {
	HistoryStore HistoryStorage
}

func (h *HistoryService) DeleteTransaction(ctx context.Context, UserID, TransactID int) error {
	return h.HistoryStore.DeleteTransaction(ctx, TransactID)
}

func NewHistoryService(historyStore HistoryStorage) *HistoryService {
	return &HistoryService{HistoryStore: historyStore}
}

func (h *HistoryService) GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error) {
	return h.HistoryStore.GetAllTransactions(ctx, UserID)
}

func (h *HistoryService) GetOneTransaction(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error) {

	transaction, err := h.HistoryStore.GetOneTransaction(ctx, TransactID)
	transaction["user_id"] = UserID
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (h *HistoryService) EditTransaction(ctx context.Context, UserID, TransactID int, Key string) error {

	err := h.HistoryStore.EditTransaction(ctx, TransactID, Key)
	if err != nil {
		return err
	}
	return nil
}
