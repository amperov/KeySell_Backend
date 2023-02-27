package main

import (
	category3 "KeySell/internal/handlers/category"
	client2 "KeySell/internal/handlers/client"
	history3 "KeySell/internal/handlers/history"
	"KeySell/internal/handlers/products"
	seller3 "KeySell/internal/handlers/seller"
	subcategory3 "KeySell/internal/handlers/subcategory"
	category2 "KeySell/internal/services/category"
	"KeySell/internal/services/client"
	history2 "KeySell/internal/services/history"
	product2 "KeySell/internal/services/product"
	seller2 "KeySell/internal/services/seller"
	subcategory2 "KeySell/internal/services/subcategory"
	"KeySell/internal/storages/category"
	"KeySell/internal/storages/history"
	"KeySell/internal/storages/product"
	"KeySell/internal/storages/seller"
	"KeySell/internal/storages/subcategory"
	"KeySell/pkg"
	"KeySell/pkg/auth"
	"KeySell/pkg/db"
	"KeySell/pkg/digi"
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	err := Init()
	if err != nil {

		logrus.Fatal(err)
		return
	}

	rtr := httprouter.New()
	PGConfig, err := db.InitPGConfig()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	PGClient, err := db.GetPGClient(context.Background(), PGConfig)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	DigiClient := digi.NewDigiClient()
	TokenManager := auth.NewTokenManager(PGClient)
	MiddleWare := auth.NewMiddleWare(TokenManager)

	CategoryStore := category.NewCategoryStorage(PGClient)
	SubcategoryStore := subcategory.NewSubcategoryStorage(PGClient)
	ProductStore := product.NewProductStorage(PGClient)
	HistoryStore := history.NewHistoryStorage(PGClient)
	SellerStore := seller.NewSellerStorage(PGClient)

	CategoryService := category2.NewCategoryService(CategoryStore, SubcategoryStore, ProductStore)
	SubcategoryService := subcategory2.NewSubcategoryService(SubcategoryStore, ProductStore, CategoryStore)
	ProductService := product2.NewProductService(ProductStore, CategoryStore)
	SellerService := seller2.NewSellerService(SellerStore, TokenManager)
	HistoryService := history2.NewHistoryService(HistoryStore)
	ClientService := client.NewClientService(ProductStore, HistoryStore, SellerStore, DigiClient, SubcategoryStore, CategoryStore)

	CategoryHandler := category3.NewCategoryHandler(MiddleWare, CategoryService)
	CategoryHandler.Register(rtr)
	SubcategoryHandler := subcategory3.NewSubcategoryHandler(MiddleWare, SubcategoryService)
	SubcategoryHandler.Register(rtr)
	ProductHandler := products.NewProductHandler(MiddleWare, ProductService)
	ProductHandler.Register(rtr)
	HistoryHandler := history3.NewHistoryHandler(MiddleWare, HistoryService)
	HistoryHandler.Register(rtr)
	SellerHandler := seller3.NewSellerHandler(MiddleWare, SellerService)
	SellerHandler.Register(rtr)
	ClientHandlers := client2.NewClientHandlers(ClientService)
	ClientHandlers.Register(rtr)

	server := pkg.NewHTTPServer(rtr)
	err = server.Run()
	if err != nil {
		logrus.Fatal(err)
		return
	}
}

func Init() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
