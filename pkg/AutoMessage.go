package pkg

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func SendMessage(Message string, InvoiceID int, Token string) error {

	logrus.Printf("InvoiceID: %d \n Token: %s\nMessage: %s", InvoiceID, Token, Message)
	bodyBytes := []byte(fmt.Sprintf(`{"message": "%s"}`, Message))
	Body := bytes.NewReader(bodyBytes)
	Response, err := http.Post(fmt.Sprintf("https://api.digiseller.ru/api/debates/v2/?token=%s&id_i=%d", Token, InvoiceID), "application/json", Body)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(Response.Body)
	if err != nil {
		logrus.Println("Message [ERROR] ", err)
		return err
	}
	logrus.Println(string(all))

	if Response.StatusCode != 200 {
		return err
	}

	return nil
}
