package auth

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type MiddleWare struct {
	tm TokenManager
}

func NewMiddleWare(tm TokenManager) MiddleWare {
	return MiddleWare{tm: tm}
}

func (w *MiddleWare) IsAuth(handle httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		var id int
		header := request.Header.Get("Authorization")
		headerArray := strings.Split(header, " ")

		id, err := w.tm.ValidateToken(headerArray[1])
		if err != nil {
			logrus.Println(err)
			return
		}
		ctx := context.WithValue(request.Context(), "user_id", id)

		handle(writer, request.WithContext(ctx), params)
	}
}
