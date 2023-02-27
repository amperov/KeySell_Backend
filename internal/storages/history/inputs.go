package history

import "github.com/sirupsen/logrus"

type Transaction struct {
	ID          int        `json:"id,omitempty"`
	Category    string     `json:"category,omitempty"`
	Subcategory string     `json:"subcategory,omitempty"`
	UniqueCode  UniqueCode `json:"unique_code,omitempty"`
	ClientEmail string     `json:"client_email,omitempty"`
	DateCheck   string     `json:"date_check,omitempty"`
	Content     string     `json:"content,omitempty"`
	Profit      float64    `json:"profit,omitempty"`
	Amount      int        `json:"amount,omitempty"`
	AmountUSD   int        `json:"amount_usd,omitempty"`
	CountGoods  int        `json:"count_goods,omitempty"`
	UniqueInv   int        `json:"unique_inv"`
	UserID      int        `json:"user_id"`
}

type UniqueCode struct {
	UniqueCode    string `json:"unique_code,omitempty"`
	State         string `json:"state,omitempty"`
	DateCheck     string `json:"date_check,omitempty"`
	DateDelivery  string `json:"date_delivery,omitempty"`
	DateConfirmed string `json:"date_confirmed,omitempty"`
}

func (t *Transaction) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["category_name"] = t.Category
	m["subcategory_name"] = t.Subcategory
	m["client_email"] = t.ClientEmail
	m["content_key"] = t.Content
	m["amount"] = t.Amount
	m["profit"] = t.Profit
	m["amount_usd"] = t.AmountUSD
	m["count"] = t.CountGoods
	m["unique_inv"] = t.UniqueInv
	m["user_id"] = t.UserID

	m["unique_code"] = t.UniqueCode.UniqueCode
	m["date_check"] = t.UniqueCode.DateCheck
	m["date_delivery"] = t.UniqueCode.DateDelivery
	m["date_confirmed"] = t.UniqueCode.DateConfirmed
	m["state"] = t.UniqueCode.State

	logrus.Debugf("Transaction info: %v", m)

	return m
}

type AllTransactInput struct {
	ID          int
	UniqueCode  string
	AmountUSD   int
	ClientEmail string
	Category    string
	Subcategory string
	Content     string
	DateCheck   string
	State       string
}

func (i *AllTransactInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = i.ID
	m["unique_code"] = i.UniqueCode
	m["amount_usd"] = i.AmountUSD
	m["client_email"] = i.ClientEmail
	m["category_name"] = i.Category
	m["subcategory_name"] = i.Subcategory
	m["content_key"] = i.Content
	m["date_check"] = i.DateCheck
	m["state"] = i.State
	return m
}
