package digi

type DigiInput struct {
	Retval          int     `json:"retval,omitempty"`
	Retdesc         string  `json:"retdesc,omitempty"`
	Inv             int     `json:"inv,omitempty"`
	IdGoods         int     `json:"id_goods,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
	TypeCurr        string  `json:"type_curr,omitempty"`
	Profit          float64 `json:"profit,omitempty"`
	AmountUsd       float64 `json:"amount_usd,omitempty"`
	DatePay         string  `json:"date_pay,omitempty"`
	Email           string  `json:"email,omitempty"`
	AgentId         int     `json:"agent_id,omitempty"`
	AgentPercent    float64 `json:"agent_percent,omitempty"`
	UnitGoods       int     `json:"unit_goods,omitempty"`
	CntGoods        float64 `json:"cnt_goods,omitempty"`
	PromoCode       string  `json:"promo_code,omitempty"`
	BonusCode       string  `json:"bonus_code,omitempty"`
	CartUid         string  `json:"cart_uid,omitempty"`
	UniqueCodeState struct {
		State         int    `json:"state,omitempty"`
		DateCheck     string `json:"date_check,omitempty"`
		DateDelivery  string `json:"date_delivery,omitempty"`
		DateConfirmed string `json:"date_confirmed,omitempty"`
		DateRefuted   string `json:"date_refuted,omitempty"`
	} `json:"unique_code_state,omitempty"`
	Options []struct {
		Id    int    `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"options,omitempty"`
}
