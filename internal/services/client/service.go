package client

import (
	"KeySell/pkg"
	"KeySell/pkg/tooling"
	"context"
	"github.com/sirupsen/logrus"
	"log"
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
	GetIDBySubItem(ctx context.Context, SubItemID int) (int, error)
	GetCatIDByID(ctx context.Context, SubItemID int) (int, error)
	IsComposite(ctx context.Context, SubCatID int) (bool, error)
	PrecheckIsComposite(ctx context.Context, SubCatID int) (bool, error)
}

type CategoryStore interface {
	GetIDByTitle(ctx context.Context, Title string) (int, string, error)
}

type ProductStorage interface {
	DeleteOne(ctx context.Context, ProdID int) error
	GetForClient(ctx context.Context, SubcategoryID, Count int) ([]map[string]interface{}, error)
	PreCheck(ctx context.Context, SubItemID int) bool
}
type ClientService struct {
	ProdStore     ProductStorage
	HistoryStore  HistoryStorage
	SellerStore   SellerStorage
	Digi          DigiClient
	SubcatStore   SubcategoryStore
	CategoryStore CategoryStore
	Select        *tooling.Tool
}

func NewClientService(prodStore ProductStorage, historyStore HistoryStorage, sellerStore SellerStorage, digi DigiClient, subcatStore SubcategoryStore, categoryStore CategoryStore, tool *tooling.Tool) *ClientService {
	return &ClientService{ProdStore: prodStore, HistoryStore: historyStore, SellerStore: sellerStore, Digi: digi, SubcatStore: subcatStore, CategoryStore: categoryStore, Select: tool}
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
		logrus.Println("Getting from history [DENIED]")
		// 2 Case: User not get products before
		Token, err := c.Digi.Auth(ctx, SellerID, SellerKey)
		if err != nil {
			logrus.Println("DigiAuth: ", err)
			return nil, err
		}

		Count, CategoryTitle, SubcategoryTitle, MapForHistory, err := c.Digi.GetInfo(ctx, UniqueCode, Token)
		if err != nil {
			logrus.Println("GetInfo", err)
			return nil, err
		}

		CategoryID, Message, err := c.CategoryStore.GetIDByTitle(ctx, CategoryTitle)
		if err != nil {
			logrus.Println("Get CatID By TITLE: ", err)
			return nil, err
		}

		SubcategoryID, err := c.SubcatStore.GetIDByTitle(ctx, SubcategoryTitle, CategoryID)
		if err != nil {
			logrus.Println("Get SubCatID By TITLE: ", err)
			return nil, err
		}

		var ProdFromStore []map[string]interface{}
		composite, err := c.SubcatStore.IsComposite(ctx, SubcategoryID)
		if err != nil {
			logrus.Printf("Is Composite: %v", err)
			return nil, err
		}
		logrus.Println("Composite: ", composite)
		if composite {
			ProdFromStore, err = c.Select.SelectTool(ctx, SubcategoryID, CategoryID)
			if err != nil {
				log.Println("SelectTool ", err)
				return nil, err
			}
		} else {
			ProdFromStore, err = c.ProdStore.GetForClient(ctx, SubcategoryID, Count)
			if err != nil {
				logrus.Println(err)
				return nil, err
			}
		}

		var content string
		for _, prod := range ProdFromStore {
			prod["client_email"] = MapForHistory["client_email"]
			prod["unique_code"] = MapForHistory["unique_code"]
			prod["unique_inv"] = MapForHistory["unique_inv"]
			content += prod["content_key"].(string) + "\n"
		}

		pkg.SendMessage(Message, MapForHistory["unique_inv"].(int), Token)

		MapForHistory["created_at"] = ProdFromStore[0]["created_at"]
		MapForHistory["content_key"] = content

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

		var NewMap []map[string]interface{}
		NewMap = append(NewMap, MapForHistory)
		return NewMap, nil
	}

	return ProdFromTx, nil
}

func (c *ClientService) Check(ctx context.Context, SubItemID int) bool {
	ID, err := c.SubcatStore.GetIDBySubItem(ctx, SubItemID)
	if err != nil {
		logrus.Printf("Get ID By SubItemID: %v", err)
		return false
	}

	IsComposite, err := c.SubcatStore.PrecheckIsComposite(ctx, SubItemID)
	if err != nil {
		logrus.Printf("Checking for Composite: %v", err)
		return false
	}
	logrus.Println("Is Composite: ", IsComposite)
	CatID, err := c.SubcatStore.GetCatIDByID(ctx, ID)
	if err != nil {
		logrus.Printf("Get CatID by SubCatID: %v", err)
		return false
	}

	var check bool
	if IsComposite {
		check, err = c.Select.SelectToolCheck(ctx, ID, CatID)
		if err != nil {
			logrus.Printf("Select Tool Error: %v", err)
			return false
		}
		return check
	} else if !IsComposite {
		check = c.ProdStore.PreCheck(ctx, SubItemID)
	}
	return check
}
