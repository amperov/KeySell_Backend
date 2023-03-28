package seller

import (
	"KeySell/pkg/auth"
	"KeySell/pkg/gen"
	"KeySell/pkg/mailer"
	"KeySell/pkg/tooling"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type Service interface {
	SignUp(ctx context.Context, m map[string]interface{}) (int, error)
	SignIn(ctx context.Context, m map[string]interface{}) (string, error)
	UpdateData(ctx context.Context, m map[string]interface{}, UserID int) error
	GetInfo(ctx context.Context, UserID int) (map[string]interface{}, error)
}

type SellerHandler struct {
	ware auth.MiddleWare
	g    gen.Gen
	s    Service
	m    mailer.Mailer
}

func NewSellerHandler(ware auth.MiddleWare, s Service, mailer2 mailer.Mailer, gen2 gen.Gen) *SellerHandler {
	return &SellerHandler{ware: ware, s: s, g: gen2, m: mailer2}
}

func (s *SellerHandler) Register(r *httprouter.Router) {
	r.POST("/api/auth/sign-in", s.AuthUser)
	r.POST("/api/auth/sign-up", s.CreateUser)
	r.PATCH("/api/seller/update", s.ware.IsAuth(s.UpdateData))
	r.GET("/api/seller/info", s.ware.IsAuth(s.GetInfo))
	r.GET("/api/seller/me", s.ware.IsAuth(s.IsAuth))
	r.GET("/api/seller/recover", s.ware.IsAuth(s.RecoverPassword))
}

func (s *SellerHandler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input SignUpInput

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	logrus.Println(input)
	id, err := s.s.SignUp(r.Context(), input.ToMap())

	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"success": "user with ID %d created"}`, id)))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(201)
}

func (s *SellerHandler) AuthUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input SignInInput

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	token, err := s.s.SignIn(r.Context(), input.ToMap())
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"JWT": "%s"}`, token)))
}

func (s *SellerHandler) UpdateData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var upd UpdateInput
	UserID := r.Context().Value("user_id").(int)

	err := tooling.GetFromBody(r.Body, &upd)
	if err != nil {
		return
	}

	err = s.s.UpdateData(r.Context(), upd.ToMap(), UserID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

func (s *SellerHandler) GetInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	info, err := s.s.GetInfo(r.Context(), UserID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		return
	}

	marshal, err := json.Marshal(info)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(marshal)
}

func (s *SellerHandler) IsAuth(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	UserID := r.Context().Value("user_id").(int)
	if UserID == 0 {
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errors.New("invalid token"))))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(200)
}

func (s *SellerHandler) RecoverPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	UserID := r.Context().Value("user_id").(int)

	info, err := s.s.GetInfo(r.Context(), UserID)
	if err != nil {
		logrus.Errorln("Get Info: ", err)
		return
	}

	if info["email"].(string) == "" {
		w.Write([]byte(`{"error": "email не указан"}`))
		return
	}

	NewPassword := s.g.Generate()
	var UpdateInput UpdateInput
	UpdateInput.Password = NewPassword

	err = s.s.UpdateData(r.Context(), UpdateInput.ToMap(), UserID)
	if err != nil {
		logrus.Errorln("Update Password: ", err)
		return
	}

	s.m.SendNewPassword(info["email"].(string), NewPassword)
}
