package pkg

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(router *httprouter.Router) *HTTPServer {
	return &HTTPServer{Server: &http.Server{Handler: router}}
}

func (s *HTTPServer) Run() error {
	s.Server.Addr = ":" + viper.GetString("app.port")

	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
