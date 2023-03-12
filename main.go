package main

/*
import (
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
)

type Resp2 struct {
	Retval  int         `json:"retval"`
	Retdesc interface{} `json:"retdesc"`
	Errors  interface{} `json:"errors"`
	Content struct {
		ItemId       int         `json:"item_id"`
		CartUid      interface{} `json:"cart_uid"`
		Name         string      `json:"name"`
		Amount       float64     `json:"amount"`
		CurrencyType string      `json:"currency_type"`
		InvoiceState int         `json:"invoice_state"`
		PurchaseDate string      `json:"purchase_date"`
		AgentId      int         `json:"agent_id"`
		AgentPercent float64     `json:"agent_percent"`
		QueryString  interface{} `json:"query_string"`
		UnitGoods    float64     `json:"unit_goods"`
		CntGoods     float64     `json:"cnt_goods"`
		PromoCode    interface{} `json:"promo_code"`
		BonusCode    interface{} `json:"bonus_code"`
		Feedback     struct {
			Deleted      bool   `json:"deleted"`
			Feedback     string `json:"feedback"`
			FeedbackType string `json:"feedback_type"`
			Comment      string `json:"comment"`
		} `json:"feedback"`
		UniqueCodeState struct {
			State         int         `json:"state"`
			DateCheck     interface{} `json:"date_check"`
			DateDelivery  interface{} `json:"date_delivery"`
			DateRefuted   interface{} `json:"date_refuted"`
			DateConfirmed interface{} `json:"date_confirmed"`
		} `json:"unique_code_state"`
		Options []struct {
			Id         int    `json:"id"`
			Name       string `json:"name"`
			UserData   string `json:"user_data"`
			UserDataId *int   `json:"user_data_id"`
		} `json:"options"`
		BuyerInfo struct {
			PaymentMethod string      `json:"payment_method"`
			Account       interface{} `json:"account"`
			Email         string      `json:"email"`
			Phone         interface{} `json:"phone"`
			Skype         interface{} `json:"skype"`
			Whatsapp      interface{} `json:"whatsapp"`
			IpAddress     interface{} `json:"ip_address"`
		} `json:"buyer_info"`
		Owner     int    `json:"owner"`
		DayLock   int    `json:"day_lock"`
		LockState string `json:"lock_state"`
	} `json:"content"`
}
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

type Body struct {
	SellerID  int    `json:"seller_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
}

type Resp struct {
	Retval   int    `json:"retval,omitempty"`
	Token    string `json:"token,omitempty"`
	SellerID int    `json:"seller_id,omitempty"`
	Valid    string `json:"valid_thru,omitempty"`
}

func main() {
	/*var body Body
	var respStr Resp
	SellerKey := "6976D6752D0B47E492EFC7E75A8BCD42"
	SellerID := 1123864
	UniqueCode := "12C678EEFB334B6C"

	body.SellerID = SellerID
	body.Timestamp = time.Now().Unix()

	hash := sha256.New()
	hash.Write([]byte(SellerKey + strconv.Itoa(int(body.Timestamp))))

	body.Sign = hex.EncodeToString(hash.Sum(nil))

	BodyMarshalled, err := json.Marshal(body)
	if err != nil {
		logrus.Debugf("marshalling error: %v", err)
		return
	}
	reader := bytes.NewReader(BodyMarshalled)

	resp, err := http.Post("https://api.digiseller.ru/api/apilogin", "application/json", reader)
	if err != nil {
		logrus.Debugf("http.Post error: %v", err)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Debugf("ReadAll error: %v", err)
		return
	}

	err = json.Unmarshal(respBody, &respStr)
	if err != nil {
		logrus.Debugf("unmarshal error: %v", err)
		return
	}
	logrus.Println(respStr)
	logrus.Println("Token: ", respStr.Token)

	/////////////////////////////

	log.Println("Get Products")
	var input DigiInput

	resp, err = http.Get(fmt.Sprintf("https://api.digiseller.ru/api/purchases/unique-code/%s?token=%s", UniqueCode, respStr.Token))
	if err != nil {
		log.Println(err)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		log.Println(err)
		return
	}
	logrus.Printf("DigiInput: %+v", input)

	//////////////////////////////////////

	var Input Resp2
	get, err := http.Get(fmt.Sprintf("https://api.digiseller.ru/api/purchase/info/%d?token=%s", input.Inv, respStr.Token))
	if err != nil {
		log.Println(err)
		return
	}

	All2, err := io.ReadAll(get.Body)
	if err != nil {
		logrus.Println(err)
		return
	}

	err = json.Unmarshal(All2, &Input)
	if err != nil {
		logrus.Println(err)
		return
	}
	fmt.Printf("Second Input: %+v", Input)
	s := reflect.Struct.String()
	logrus.Println(s)
	logrus.Println(govalidator.IsEmail("Aoaoaoao"))
	logrus.Println(govalidator.IsExistingEmail("unflat.gopher@gl.com"))
}
*/
