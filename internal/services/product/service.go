package product

import (
	"context"
	"errors"
)

type ProductStorage interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, ProdID int) (int, error)
	Delete(ctx context.Context, ProdID int) error
}
type CatStore interface {
	Belong(ctx context.Context, UserID, CatID int) bool
}

type ProductService struct {
	ProdStore ProductStorage
	CatStore  CatStore
}

func NewProductService(prodStore ProductStorage, catStore CatStore) *ProductService {
	return &ProductService{ProdStore: prodStore, CatStore: catStore}
}

func (p *ProductService) Create(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error) {
	belong := p.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return 0, errors.New("you dont have permissions")
	}

	ID, err := p.ProdStore.Create(ctx, m)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (p *ProductService) Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID, ProdID int) (int, error) {
	belong := p.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return 0, errors.New("you dont have permissions")
	}
	ID, err := p.ProdStore.Update(ctx, m, ProdID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (p *ProductService) Delete(ctx context.Context, UserID, CatID, SubCatID, ProdID int) error {
	belong := p.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return errors.New("you dont have permissions")
	}
	err := p.ProdStore.Delete(ctx, ProdID)
	if err != nil {
		return err
	}
	return nil
}
