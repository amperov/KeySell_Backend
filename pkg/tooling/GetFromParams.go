package tooling

import (
	"github.com/julienschmidt/httprouter"
	"strconv"
)

func GetCategoryID(params httprouter.Params) int {
	catID := params.ByName("cat_id")
	CatID, err := strconv.Atoi(catID)
	if err != nil {
		return 0
	}
	return CatID
}

func GetSubcategoryID(params httprouter.Params) int {
	subcatID := params.ByName("subcat_id")
	SubcatID, err := strconv.Atoi(subcatID)
	if err != nil {
		return 0
	}
	return SubcatID
}

func GetProductID(params httprouter.Params) int {
	prodID := params.ByName("prod_id")
	ProdID, err := strconv.Atoi(prodID)
	if err != nil {
		return 0
	}
	return ProdID
}

func GetTwoIDs(params httprouter.Params) (int, int) {
	CatID := GetCategoryID(params)
	SubCatID := GetSubcategoryID(params)
	return CatID, SubCatID
}

func GetAllIDs(params httprouter.Params) (int, int, int) {
	CatID := GetCategoryID(params)
	SubCatID := GetSubcategoryID(params)
	ProdID := GetProductID(params)
	return CatID, SubCatID, ProdID
}
func GetTranID(params httprouter.Params) int {
	tranID := params.ByName("tran_id")
	TranID, err := strconv.Atoi(tranID)
	if err != nil {
		return 0
	}
	return TranID
}
