package seller

type Seller struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	SellerID  int    `json:"seller_id,omitempty"`
	SellerKey string `json:"seller_key,omitempty"`
}

func (m *Seller) ToMap() map[string]interface{} {

	var SellerMap = make(map[string]interface{})

	SellerMap["username"] = m.Username

	SellerMap["email"] = m.Email

	SellerMap["seller_id"] = m.SellerID

	SellerMap["seller_key"] = m.SellerKey

	return SellerMap
}
