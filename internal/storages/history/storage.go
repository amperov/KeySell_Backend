package history

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
)

var table = "transactions"

type HistoryStorage struct {
	c *pgxpool.Pool
}

func (h *HistoryStorage) EditTransaction(ctx context.Context, TransactID int, Key string) error {
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": TransactID}).Set("content_key", Key).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = h.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewHistoryStorage(c *pgxpool.Pool) *HistoryStorage {
	return &HistoryStorage{c: c}
}

func (h *HistoryStorage) GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error) {
	var transacts AllTransactInput
	m := make(map[string]interface{})

	query, args, err := squirrel.
		Select("id", "category_name", "subcategory_name", "unique_code",
			"content_key", "state", "amount", "date_check", "client_email", "unique_inv", "created_at").
		From(table).Where(squirrel.Eq{"user_id": UserID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := h.c.Query(ctx, query, args...)
	if err != nil {
		logrus.Debugf("error Query: %v", err)
		return nil, err
	}

	var arrayMap []map[string]interface{}

	for rows.Next() {
		err := rows.Scan(&transacts.ID, &transacts.Category, &transacts.Subcategory, &transacts.UniqueCode,
			&transacts.Content, &transacts.State, &transacts.AmountUSD, &transacts.DateCheck, &transacts.ClientEmail, &transacts.UniqueInv, &transacts.CreatedAt)
		if err != nil {
			logrus.Debugf("error scanning: %v", err)
			return nil, err
		}
		arrayMap = append(arrayMap, transacts.ToMap())
	}
	m["transactions"] = arrayMap

	return m, nil
}

func (h *HistoryStorage) GetOneTransaction(ctx context.Context, TransactID int) (map[string]interface{}, error) {
	var transtact Transaction
	query, args, err := squirrel.
		Select("category_name", "subcategory_name", "unique_code",
			"client_email", "amount", "profit", "count",
			"unique_inv", "date_delivery", "date_confirmed",
			"content_key", "state", "amount_usd", "date_check").
		From(table).Where(squirrel.Eq{"id": TransactID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := h.c.QueryRow(ctx, query, args...)
	err = row.Scan(&transtact.Category, &transtact.Subcategory, &transtact.UniqueCode.UniqueCode,
		&transtact.ClientEmail, &transtact.Amount, &transtact.Profit, &transtact.CountGoods,
		&transtact.UniqueInv, &transtact.UniqueCode.DateDelivery, &transtact.UniqueCode.DateConfirmed,
		&transtact.Content, &transtact.UniqueCode.State, &transtact.AmountUSD, &transtact.UniqueCode.DateCheck)
	if err != nil {
		logrus.Debugf("error scanning: %v", err)
		return nil, err
	}

	transtact.ID = TransactID
	var mm = make(map[string]interface{})
	mm["transaction"] = transtact.ToMap()
	return mm, nil
}

func (h *HistoryStorage) GetByUC(ctx context.Context, UniqueCode string, UserID int) ([]map[string]interface{}, error) {
	var transtact Transaction
	query, args, err := squirrel.
		Select("id", "category_name", "subcategory_name",
			"client_email", "amount", "profit", "count",
			"unique_inv", "date_delivery", "date_confirmed",
			"content_key", "state", "amount_usd", "date_check").
		From(table).Where(squirrel.Eq{"unique_code": UniqueCode, "user_id": UserID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := h.c.QueryRow(ctx, query, args...)

	err = row.Scan(&transtact.ID, &transtact.Category, &transtact.Subcategory, &transtact.ClientEmail,
		&transtact.Amount, &transtact.Profit, &transtact.CountGoods, &transtact.UniqueInv, &transtact.UniqueCode.DateDelivery, &transtact.UniqueCode.DateConfirmed, &transtact.Content,
		&transtact.UniqueCode.State, &transtact.AmountUSD, &transtact.DateCheck)
	if err != nil {
		logrus.Debugf("error scanning: %v", err)
		return nil, err
	}
	transtact.UniqueCode.UniqueCode = UniqueCode
	var ArrayTransactions []map[string]interface{}
	ArrayTransactions = append(ArrayTransactions, transtact.ToMap())

	return ArrayTransactions, nil
}

func (h *HistoryStorage) SetTransaction(ctx context.Context, m map[string]interface{}, UserID int) error {
	m["user_id"] = UserID

	query, args, err := squirrel.Insert(table).SetMap(m).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = h.c.Exec(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
