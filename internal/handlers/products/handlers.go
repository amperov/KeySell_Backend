package products

import (
	"KeySell/pkg/auth"
	"KeySell/pkg/tooling"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
)

type ProductService interface {
	Create(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID, ProdID int) (int, error)
	Delete(ctx context.Context, UserID, CatID, SubCatID, ProdID int) error
}
type ProductHandler struct {
	ware auth.MiddleWare
	ps   ProductService
}

func NewProductHandler(ware auth.MiddleWare, ps ProductService) *ProductHandler {
	return &ProductHandler{ware: ware, ps: ps}
}

func (h *ProductHandler) Register(r *httprouter.Router) {
	r.POST("/api/seller/category/:cat_id/subcategory/:subcat_id/one", h.ware.IsAuth(h.CreateOne))
	r.POST("/api/seller/category/:cat_id/subcategory/:subcat_id/many", h.ware.IsAuth(h.CreateMany))
	r.PATCH("/api/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.UpdateProduct))
	r.DELETE("/api/seller/category/:cat_id/subcategory/:subcat_id/products/:product_id", h.ware.IsAuth(h.DeleteProduct))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input UpdateProductInput

	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID, ProductID := tooling.GetAllIDs(params)

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	input.SubCatID = SubCatID
	id, err := h.ps.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(400)
		logrus.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"status": "product %d updated"}`, id)))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID, ProductID := tooling.GetAllIDs(params)

	err := h.ps.Delete(r.Context(), UserID, CatID, SubCatID, ProductID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte(`{"success" : "product deleted"}`))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
}

func (h *ProductHandler) CreateOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input CreateProductInput

	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID := tooling.GetTwoIDs(params)

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}
	input.SubCatID = SubCatID

	id, err := h.ps.Create(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		log.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"status": "product %d created"}`, id)))
}

func (h *ProductHandler) CreateMany(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input CreateProductInput

	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID := tooling.GetTwoIDs(params)

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	Products := strings.Split(input.Content, "\n")

	count := 0
	for i := 0; i < len(Products); i++ {
		input.Content = Products[i]
		input.SubCatID = SubCatID
		_, err := h.ps.Create(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
		if err != nil {
			log.Println(err)
			continue
		}
		count++
	}
}
