package seller

type Seller struct {
	Username  string `json:"username,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Password  string `json:"password,omitempty"`
	SellerID  int    `json:"seller_id,omitempty"`
	SellerKey string `json:"seller_key,omitempty"`
}

func (m *Seller) ToMap() map[string]interface{} {

	var SellerMap = make(map[string]interface{})

	SellerMap["firstname"] = m.Firstname

	SellerMap["username"] = m.Username

	SellerMap["lastname"] = m.Lastname

	SellerMap["seller_id"] = m.SellerID

	SellerMap["seller_key"] = m.SellerKey

	return SellerMap
}
