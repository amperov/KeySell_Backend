package category

import (
	"KeySell/pkg/auth"
	"KeySell/pkg/tooling"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

//Replace interfaces to Structures

type CatService interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error)
	Delete(ctx context.Context, UserID int, CatID int) error
	GetOne(ctx context.Context, UserID int, CatID int) (map[string]interface{}, error)
	GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error)
}

type CategoryHandler struct {
	ware auth.MiddleWare
	cat  CatService
}

func NewCategoryHandler(ware auth.MiddleWare, cat CatService) *CategoryHandler {
	return &CategoryHandler{ware: ware, cat: cat}
}

func (h *CategoryHandler) Register(r *httprouter.Router) {
	r.GET("/api/seller/category", h.ware.IsAuth(h.GetAllCategory))
	r.GET("/api/seller/category/:cat_id", h.ware.IsAuth(h.GetCategory))
	r.POST("/api/seller/category", h.ware.IsAuth(h.CreateCategory))
	r.PATCH("/api/seller/category/:cat_id", h.ware.IsAuth(h.UpdateCategory))
	r.DELETE("/api/seller/category/:cat_id", h.ware.IsAuth(h.DeleteCategory))
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input CreateCategoryInput
	UserID := r.Context().Value("user_id").(int)

	if UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}
	input.UserID = UserID

	id, err := h.cat.Create(r.Context(), input.ToMap())

	_, err = w.Write([]byte(fmt.Sprintf(`{"success" : "category with ID %d created"}`, id)))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var upd UpdateCategoryInput
	UserID := r.Context().Value("user_id").(int)

	CatID := tooling.GetCategoryID(params)

	err := tooling.GetFromBody(r.Body, upd)
	if err != nil {
		return
	}

	CatID, err = h.cat.Update(r.Context(), upd.ToMap(), UserID, CatID)
	if err != nil {
		return
	}

}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	CatID := tooling.GetCategoryID(params)

	err := h.cat.Delete(r.Context(), UserID, CatID)
	if err != nil {
		return
	}
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	CatID := tooling.GetCategoryID(params)

	logrus.Println("CatID From GetCategoryID:", CatID)
	Category, err := h.cat.GetOne(r.Context(), UserID, CatID)
	if err != nil {
		logrus.Println(err)
		return
	}

	MapCategory, err := json.Marshal(Category)
	if err != nil {
		logrus.Println(err)
		return
	}
	w.Write(MapCategory)
}

func (h *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	UserID := r.Context().Value("user_id").(int)

	Categories, err := h.cat.GetAll(r.Context(), UserID)
	if err != nil {
		logrus.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
	MapCategories, err := json.Marshal(Categories)
	if err != nil {
		logrus.Println(err)
		return
	}
	w.Write(MapCategories)
}
