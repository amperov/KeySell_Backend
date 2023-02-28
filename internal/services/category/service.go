package category

import (
	"context"
	"errors"
)

type CategoryStorage interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, CatID int) (int, error)
	Delete(ctx context.Context, CatID int) error
	GetOne(ctx context.Context, CatID int) (map[string]interface{}, error)
	GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error)
	Belong(ctx context.Context, UserID, CategoryID int) bool
}

type SubcategoryStorage interface {
	GetAll(ctx context.Context, CatID int) ([]map[string]interface{}, error)
	GetCount(ctx context.Context, CatID int) (int, error)
}

type ProductStorage interface {
	GetCount(ctx context.Context, CatID int) (int, error)
}

type CategoryService struct {
	CatStore    CategoryStorage
	SubCatStore SubcategoryStorage
	ProdStore   ProductStorage
}

func NewCategoryService(catStore CategoryStorage, subCatStore SubcategoryStorage, prodStore ProductStorage) *CategoryService {
	return &CategoryService{CatStore: catStore, SubCatStore: subCatStore, ProdStore: prodStore}
}

func (c *CategoryService) Create(ctx context.Context, m map[string]interface{}) (int, error) {

	return c.CatStore.Create(ctx, m)
}

func (c *CategoryService) Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error) {
	belong := c.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return 0, errors.New("you dont have permissions")
	}

	return c.CatStore.Update(ctx, m, CatID)
}

func (c *CategoryService) Delete(ctx context.Context, UserID int, CatID int) error {
	belong := c.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return errors.New("you dont have permissions")
	}
	return c.CatStore.Delete(ctx, CatID)
}

func (c *CategoryService) GetOne(ctx context.Context, UserID int, CatID int) (map[string]interface{}, error) {
	belong := c.CatStore.Belong(ctx, UserID, CatID)
	if belong == false {
		return nil, errors.New("you dont have permissions")
	}

	Category, err := c.CatStore.GetOne(ctx, CatID)
	if err != nil {
		return nil, err
	}
	SubCategories, err := c.SubCatStore.GetAll(ctx, CatID)
	if err != nil {
		return nil, err
	}

	for _, subcat := range SubCategories {
		count, err := c.ProdStore.GetCount(ctx, subcat["id"].(int))
		if err != nil {
			return nil, err
		}
		subcat["count_products"] = count
	}
	Category["user_id"] = UserID
	Category["subcategories"] = SubCategories
	return Category, nil
}

func (c *CategoryService) GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error) {
	all, err := c.CatStore.GetAll(ctx, UserID)
	if err != nil {
		return nil, err
	}

	for _, m := range all {
		count, err := c.SubCatStore.GetCount(ctx, m["id"].(int))
		if err != nil {
			count = 0
		}
		m["user_id"] = UserID
		m["count_subcategories"] = count
	}
	return all, nil
}
