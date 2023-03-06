package digi

import (
	"KeySell/internal/storages/history"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

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

type DigiClient struct {
}

func NewDigiClient() *DigiClient {
	return &DigiClient{}
}

func (c *DigiClient) Auth(ctx context.Context, SellerID int, SellerKey string) (string, error) {
	var body Body

	var respStr Resp
	//Searching SellerID by Username or integrate Seller Storage

	body.SellerID = SellerID
	body.Timestamp = time.Now().Unix()

	hash := sha256.New()
	hash.Write([]byte(SellerKey + strconv.Itoa(int(body.Timestamp))))
	body.Sign = hex.EncodeToString(hash.Sum(nil))

	BodyMarshalled, err := json.Marshal(body)
	if err != nil {
		logrus.Debugf("marshalling error: %v", err)
		return "", err
	}
	reader := bytes.NewReader(BodyMarshalled)

	resp, err := http.Post("https://api.digiseller.ru/api/apilogin", "application/json", reader)
	if err != nil {
		logrus.Debugf("http.Post error: %v", err)
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Debugf("ReadAll error: %v", err)
		return "", err
	}

	err = json.Unmarshal(respBody, &respStr)
	if err != nil {
		logrus.Debugf("unmarshal error: %v", err)
		return "", err
	}
	return respStr.Token, nil
}

func (c *DigiClient) GetInfo(ctx context.Context, UniqueCode, Token string) (int, string, string, map[string]interface{}, error) {
	log.Println("Get Products")
	var input DigiInput
	var tran history.Transaction

	resp, err := http.Get(fmt.Sprintf("https://api.digiseller.ru/api/purchases/unique-code/%s?token=%s", UniqueCode, Token))
	if err != nil {
		log.Println(err)
		return 0, "", "", nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return 0, "", "", nil, err
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		log.Println(err)
		return 0, "", "", nil, err
	}

	tran.UniqueInv = input.Inv
	tran.UniqueCode.UniqueCode = UniqueCode
	tran.UniqueCode.DateConfirmed = input.UniqueCodeState.DateConfirmed
	tran.UniqueCode.DateDelivery = input.UniqueCodeState.DateDelivery
	tran.UniqueCode.DateCheck = input.UniqueCodeState.DateCheck
	tran.CountGoods = int(input.CntGoods)
	tran.Amount = int(input.Amount)
	tran.AmountUSD = int(input.AmountUsd)
	tran.Category = input.Options[0].Name
	tran.Subcategory = input.Options[0].Value
	tran.ClientEmail = input.Email
	tran.Profit = input.Profit

	return tran.CountGoods, tran.Category, tran.Subcategory, tran.ToMap(), nil
}
