package client

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type ClientService interface {
	Get(ctx context.Context, UniqueCode string, Username string) ([]map[string]interface{}, error)
	Check(ctx context.Context, ItemID int) (bool, error)
}

type ClientHandlers struct {
	c ClientService
}

func NewClientHandlers(c ClientService) *ClientHandlers {
	return &ClientHandlers{c: c}
}

func (h *ClientHandlers) Register(r *httprouter.Router) {
	r.GET("/api/client/:username", h.GetProducts)
	r.POST("/api/precheck", h.PreCheck)

}

func (h *ClientHandlers) GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uniqueCode := r.URL.Query().Get("uniquecode")
	name := params.ByName("username")

	prods, err := h.c.Get(r.Context(), uniqueCode, name)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	prodsMarshalled, err := json.Marshal(prods)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}

	w.Write(prodsMarshalled)

}

// Precheck
type Request struct {
	XMLName xml.Name `xml:"request"`
	Text    string   `xml:",chardata"`
	Product struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id"`
		Cnt  string `xml:"cnt"`
		Lang string `xml:"lang"`
	} `xml:"product"`
	Options struct {
		Text   string `xml:",chardata"`
		Option []struct {
			Text int    `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Type string `xml:"type,attr"`
		} `xml:"option"`
	} `xml:"options"`
}

func (h *ClientHandlers) PreCheck(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var input Request

	all, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(all, &input)
	if err != nil {
		return
	}

	log.Printf("%+v", input)
	var SubItemID int

	for _, option := range input.Options.Option {
		if option.Type == "radio" {
			SubItemID = option.Text
		}

	}

	check, err := h.c.Check(request.Context(), SubItemID)
	if err != nil {
		log.Print(err)
		return
	}
	if check == false {
		w.WriteHeader(400)
		w.Write([]byte(`"error": "we haven't this products"`))
		log.Println(err)
		return
	}
	w.Write([]byte(`{"error": ""}`))
	w.WriteHeader(200)
}
