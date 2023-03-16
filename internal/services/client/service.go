package client

import (
	"KeySell/pkg"
	"context"
	"github.com/sirupsen/logrus"
)

type DigiClient interface {
	Auth(ctx context.Context, SellerID int, SellerKey string) (string, error)
	//GetInfo Return: Count, CategoryTitle, SubcategoryTitle, Error
	GetInfo(ctx context.Context, UniqueCode, Token string) (int, string, string, map[string]interface{}, error)
}
type HistoryStorage interface {
	SetTransaction(ctx context.Context, m map[string]interface{}, UserID int) error
	GetByUC(ctx context.Context, UniqueCode string, UserID int) ([]map[string]interface{}, error)
	// TODO in Seller Service EditTransaction(ctx context.Context, ProdContent string) error
}
type SellerStorage interface {
	GetIDByUsername(ctx context.Context, Username string) (int, error)
	GetDataByID(ctx context.Context, UserID int) (int, string, error)
}

type SubcategoryStore interface {
	GetIDByTitle(ctx context.Context, Title string, CategoryID int) (int, error)
}

type CategoryStore interface {
	GetIDByTitle(ctx context.Context, Title string) (int, string, error)
}

type ProductStorage interface {
	DeleteOne(ctx context.Context, ProdID int) error
	GetForClient(ctx context.Context, SubcategoryID, Count int) ([]map[string]interface{}, error)
	PreCheck(ctx context.Context, SubItemID int) (bool, error)
}
type ClientService struct {
	ProdStore     ProductStorage
	HistoryStore  HistoryStorage
	SellerStore   SellerStorage
	Digi          DigiClient
	SubcatStore   SubcategoryStore
	CategoryStore CategoryStore
}

func NewClientService(prodStore ProductStorage, historyStore HistoryStorage, sellerStore SellerStorage, digi DigiClient, subcatStore SubcategoryStore, categoryStore CategoryStore) *ClientService {
	return &ClientService{ProdStore: prodStore, HistoryStore: historyStore, SellerStore: sellerStore, Digi: digi, SubcatStore: subcatStore, CategoryStore: categoryStore}
}

func (c *ClientService) Get(ctx context.Context, UniqueCode string, Username string) ([]map[string]interface{}, error) {
	//Getting UserID By Username
	UserID, err := c.SellerStore.GetIDByUsername(ctx, Username)
	if err != nil {
		return nil, err
	}
	// Getting Data for DigiSeller
	SellerID, SellerKey, err := c.SellerStore.GetDataByID(ctx, UserID)
	if err != nil {
		return nil, err
	}

	// 1 Case: User get products before
	ProdFromTx, err := c.HistoryStore.GetByUC(ctx, UniqueCode, UserID)
	if err != nil || ProdFromTx == nil {
		logrus.Println(err)
		// 2 Case: User not get products before
		Token, err := c.Digi.Auth(ctx, SellerID, SellerKey)
		if err != nil {
			logrus.Println("DigiAuth error: ", err)
			return nil, err
		}

		Count, CategoryTitle, SubcategoryTitle, MapForHistory, err := c.Digi.GetInfo(ctx, UniqueCode, Token)
		if err != nil {
			logrus.Println("GetInfo error", err)
			return nil, err
		}

		logrus.Println(CategoryTitle, SubcategoryTitle)

		CategoryID, Message, err := c.CategoryStore.GetIDByTitle(ctx, CategoryTitle)
		if err != nil {
			logrus.Println("Cat error: ", err)
			return nil, err
		}

		SubcategoryID, err := c.SubcatStore.GetIDByTitle(ctx, SubcategoryTitle, CategoryID)
		if err != nil {
			logrus.Println("Subcat error: ", err)
			return nil, err
		}

		ProdFromStore, err := c.ProdStore.GetForClient(ctx, SubcategoryID, Count)
		if err != nil {
			logrus.Println("Prods error: ", err)
		}

		var content string

		if ProdFromStore != nil {
			for _, prod := range ProdFromStore {
				prod["unique_code"] = MapForHistory["unique_code"]
				prod["unique_inv"] = MapForHistory["unique_inv"]

				content += prod["content_key"].(string) + "\n"
			}
			pkg.SendMessage(Message, MapForHistory["unique_inv"].(int), Token)

			MapForHistory["created_at"] = ProdFromStore[0]["created_at"]
			MapForHistory["content_key"] = content
		} else {
			MapForHistory["content_key"] = "Скоро ваш ключ будет загружен, ожидайте"
			MapForHistory["created_at"] = "Отсутствует"
			pkg.SendMessage("Данного товара сейчас нет в наличии, попробуйте обновить эту страницу "+
				"через 5-10 минут\nЕсли ключ по прежнему не отображается, свяжитесь с продавцом.\n\nThis product is currently out of stock, please try refreshing this page in 5-10 minutes\nIf the key is still not displayed, contact the seller.", MapForHistory["unique_inv"].(int), Token)
		}

		err = c.HistoryStore.SetTransaction(ctx, MapForHistory, UserID)
		if err != nil {
			logrus.Println(err)
			return nil, err
		}
		for _, Prod := range ProdFromStore {
			err := c.ProdStore.DeleteOne(ctx, Prod["id"].(int))
			if err != nil {
				return nil, err
			}
		}

		return ProdFromStore, nil
	}

	return ProdFromTx, nil
}

func (c *ClientService) Check(ctx context.Context, SubItemID int) (bool, error) {
	return c.ProdStore.PreCheck(ctx, SubItemID)
}
