package main

//
//import (
//	"bytes"
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"github.com/sirupsen/logrus"
//	"io"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//type DigiInput struct {
//	Retval          int     `json:"retval,omitempty"`
//	Retdesc         string  `json:"retdesc,omitempty"`
//	Inv             int     `json:"inv,omitempty"`
//	IdGoods         int     `json:"id_goods,omitempty"`
//	Amount          float64 `json:"amount,omitempty"`
//	TypeCurr        string  `json:"type_curr,omitempty"`
//	Profit          float64 `json:"profit,omitempty"`
//	AmountUsd       float64 `json:"amount_usd,omitempty"`
//	DatePay         string  `json:"date_pay,omitempty"`
//	Email           string  `json:"email,omitempty"`
//	AgentId         int     `json:"agent_id,omitempty"`
//	AgentPercent    float64 `json:"agent_percent,omitempty"`
//	UnitGoods       int     `json:"unit_goods,omitempty"`
//	CntGoods        float64 `json:"cnt_goods,omitempty"`
//	PromoCode       string  `json:"promo_code,omitempty"`
//	BonusCode       string  `json:"bonus_code,omitempty"`
//	CartUid         string  `json:"cart_uid,omitempty"`
//	UniqueCodeState struct {
//		State         int    `json:"state,omitempty"`
//		DateCheck     string `json:"date_check,omitempty"`
//		DateDelivery  string `json:"date_delivery,omitempty"`
//		DateConfirmed string `json:"date_confirmed,omitempty"`
//		DateRefuted   string `json:"date_refuted,omitempty"`
//	} `json:"unique_code_state,omitempty"`
//	Options []struct {
//		Id    int    `json:"id,omitempty"`
//		Name  string `json:"name,omitempty"`
//		Value string `json:"value,omitempty"`
//	} `json:"options,omitempty"`
//}
//type Body struct {
//	SellerID  int    `json:"seller_id,omitempty"`
//	Timestamp int64  `json:"timestamp,omitempty"`
//	Sign      string `json:"sign,omitempty"`
//}
//
//type Resp struct {
//	Retval   int    `json:"retval,omitempty"`
//	Token    string `json:"token,omitempty"`
//	SellerID int    `json:"seller_id,omitempty"`
//	Valid    string `json:"valid_thru,omitempty"`
//}
//
//func main() {
//	var body Body
//
//	var respStr Resp
//	//Searching SellerID by Username or integrate Seller Storage
//
//	body.SellerID = 841146
//	body.Timestamp = time.Now().Unix()
//
//	hash := sha256.New()
//	hash.Write([]byte("3A095A88BA044066B860DD4F18F84A5C" + strconv.Itoa(int(body.Timestamp))))
//	body.Sign = hex.EncodeToString(hash.Sum(nil))
//
//	BodyMarshalled, err := json.Marshal(body)
//	if err != nil {
//		logrus.Debugf("marshalling error: %v", err)
//		return
//	}
//	reader := bytes.NewReader(BodyMarshalled)
//
//	resp, err := http.Post("https://api.digiseller.ru/api/apilogin", "application/json", reader)
//	if err != nil {
//		logrus.Debugf("http.Post error: %v", err)
//		return
//	}
//
//	respBody, err := io.ReadAll(resp.Body)
//	if err != nil {
//		logrus.Debugf("ReadAll error: %v", err)
//		return
//	}
//
//	err = json.Unmarshal(respBody, &respStr)
//	if err != nil {
//		logrus.Debugf("unmarshal error: %v", err)
//		return
//	}
//	logrus.Println(respStr)
//}
//*/
