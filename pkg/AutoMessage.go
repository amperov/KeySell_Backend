package pkg

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendMessage(Message string, InvoiceID int, Token string) error {

	bodyBytes := []byte(fmt.Sprintf(`{"message": "%s"}`, Message))
	Body := bytes.NewReader(bodyBytes)
	Response, err := http.Post(fmt.Sprintf("https://api.digiseller.ru/api/debates/v2/?token=%s&id_i=%d", Token, InvoiceID), "application/json", Body)
	if err != nil {
		return err
	}

	if Response.StatusCode != 200 {
		return err
	}

	return nil
}
