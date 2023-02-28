package pkg

import (
	"github.com/spf13/viper"
	"net/http"
)

type HTTPServer struct {
	Server *http.Server
}

func NewHTTPServer(handler http.Handler) *HTTPServer {
	return &HTTPServer{Server: &http.Server{Handler: handler}}
}

func (s *HTTPServer) Run() error {
	s.Server.Addr = ":" + viper.GetString("app.port")

	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
