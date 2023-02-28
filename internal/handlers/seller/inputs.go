package seller

type SignUpInput struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	SellerID  int    `json:"seller_id,omitempty"`
	SellerKey string `json:"seller_key,omitempty"`
}

func (m SignUpInput) ToMap() map[string]interface{} {

	var SellerMap = make(map[string]interface{})

	if m.Username != "" {
		SellerMap["username"] = m.Username
	}
	if m.Password != "" {
		SellerMap["pass"] = m.Password
	}
	if m.SellerID != 0 {
		SellerMap["seller_id"] = m.SellerID
	}
	if m.SellerKey != "" {
		SellerMap["seller_key"] = m.SellerKey
	}
	return SellerMap
}

type SignInInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (m *SignInInput) ToMap() map[string]interface{} {

	var SellerMap = make(map[string]interface{})

	SellerMap["username"] = m.Username

	if m.Password != "" {
		SellerMap["password"] = m.Password
	}
	return SellerMap
}

type UpdateInput struct {
	Password  string `json:"password,omitempty"`
	SellerID  int    `json:"seller_id,omitempty"`
	SellerKey string `json:"seller_key,omitempty"`
}

func (m *UpdateInput) ToMap() map[string]interface{} {
	var SellerMap = make(map[string]interface{})

	if m.Password != "" {
		SellerMap["password"] = m.Password
	}
	if m.SellerKey != "" {
		SellerMap["seller_key"] = m.SellerKey
	}
	if m.SellerID != 0 {
		SellerMap["seller_id"] = m.SellerID
	}
	return SellerMap
}
