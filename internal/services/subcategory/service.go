package subcategory

import (
	"context"
	"errors"
)

type SubcategoryStorage interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, SubCatID int) (int, error)
	Delete(ctx context.Context, SubCatID int) error
	GetOne(ctx context.Context, SubCatID int) (map[string]interface{}, error)
}
type ProductStorage interface {
	GetAll(ctx context.Context, SubCatID int) ([]map[string]interface{}, error)
}

type CatStore interface {
	Belong(ctx context.Context, UserID, CatID int) bool
}

type SubcategoryService struct {
	SubcatStore SubcategoryStorage
	ProdStore   ProductStorage
	CatStore    CatStore
}

func NewSubcategoryService(subcatStore SubcategoryStorage, prodStore ProductStorage, catStore CatStore) *SubcategoryService {
	return &SubcategoryService{SubcatStore: subcatStore, ProdStore: prodStore, CatStore: catStore}
}

func (s *SubcategoryService) Create(ctx context.Context, m map[string]interface{}, UserID, CatID int) (int, error) {
	belong := s.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return 0, errors.New("you dont have permissions")
	}
	return s.SubcatStore.Create(ctx, m)
}

func (s *SubcategoryService) Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error) {
	belong := s.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return 0, errors.New("you dont have permissions")
	}
	return s.SubcatStore.Update(ctx, m, SubCatID)
}

func (s *SubcategoryService) Delete(ctx context.Context, UserID, CatID, SubCatID int) error {
	belong := s.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return errors.New("you dont have permissions")
	}
	return s.SubcatStore.Delete(ctx, SubCatID)
}

func (s *SubcategoryService) Get(ctx context.Context, UserID, CatID, SubCatID int) (map[string]interface{}, error) {
	belong := s.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return nil, errors.New("you dont have permissions")
	}
	SubCat, err := s.SubcatStore.GetOne(ctx, SubCatID)
	if err != nil {
		return nil, err
	}
	AllProducts, err := s.ProdStore.GetAll(ctx, SubCatID)
	if err != nil {
		return nil, err
	}
	SubCat["products"] = AllProducts

	return SubCat, nil
}
