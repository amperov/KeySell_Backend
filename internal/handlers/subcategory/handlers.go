package subcategory

import (
	"KeySell/pkg/auth"
	"KeySell/pkg/tooling"
	"github.com/sirupsen/logrus"

	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SubcatService interface {
	Create(ctx context.Context, m map[string]interface{}, UserID, CatID int) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error)
	Delete(ctx context.Context, UserID, CatID, SubCatID int) error
	Get(ctx context.Context, UserID, CatID, SubCatID int) (map[string]interface{}, error)
}

type SubcategoryHandler struct {
	ware auth.MiddleWare
	sc   SubcatService
}

func NewSubcategoryHandler(ware auth.MiddleWare, sc SubcatService) *SubcategoryHandler {
	return &SubcategoryHandler{ware: ware, sc: sc}
}

func (h *SubcategoryHandler) Register(r *httprouter.Router) {
	r.GET("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.GetSubcategory))
	r.POST("/api/seller/category/:cat_id", h.ware.IsAuth(h.CreateSubcategory))
	r.PATCH("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.UpdateSubcategory))
	r.DELETE("/api/seller/category/:cat_id/subcategory/:subcat_id", h.ware.IsAuth(h.DeleteSubcategory))
}

func (h *SubcategoryHandler) CreateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	var input CreateSubcatInput

	UserID := r.Context().Value("user_id").(int)

	CatID := tooling.GetCategoryID(params)

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	input.CatID = CatID

	id, err := h.sc.Create(r.Context(), input.ToMap(), UserID, CatID)
	if err != nil {
		w.WriteHeader(500)
		logrus.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"status": "subcategory %d created"}`, id)))

}

func (h *SubcategoryHandler) UpdateSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var input UpdateSubcatInput

	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID := tooling.GetTwoIDs(params)

	err := tooling.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}
	input.CatID = CatID

	id, err := h.sc.Update(r.Context(), input.ToMap(), UserID, CatID, SubCatID)
	if err != nil {
		w.WriteHeader(500)
		logrus.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"status": "subcategory %d updated"}`, id)))
}

func (h *SubcategoryHandler) DeleteSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID := tooling.GetTwoIDs(params)

	err := h.sc.Delete(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		logrus.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"status": "subcategory %d deleted"}`, SubCatID)))
}

func (h *SubcategoryHandler) GetSubcategory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	UserID := r.Context().Value("user_id").(int)

	CatID, SubCatID := tooling.GetTwoIDs(params)

	subcat, err := h.sc.Get(r.Context(), UserID, CatID, SubCatID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		w.WriteHeader(500)
		logrus.Println(err)
		return
	}

	bytes, err := json.Marshal(subcat)
	if err != nil {
		return
	}

	w.Write(bytes)
}
