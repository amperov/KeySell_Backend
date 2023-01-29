package category

import "context"

type CategoryStorage interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error)
	Delete(ctx context.Context, UserID int, CatID int) error
	GetOne(ctx context.Context, CatID int) (map[string]interface{}, error)
	GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error)
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

func (c *CategoryService) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	return c.CatStore.Create(ctx, m)
}

func (c *CategoryService) Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error) {
	return c.CatStore.Update(ctx, m, UserID, CatID)
}

func (c *CategoryService) Delete(ctx context.Context, UserID int, CatID int) error {
	return c.CatStore.Delete(ctx, UserID, CatID)
}

func (c *CategoryService) GetOne(ctx context.Context, UserID int, CatID int) (map[string]interface{}, error) {
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
			return nil, err
		}
		m["count_subcategories"] = count
	}
	return all, nil
}
