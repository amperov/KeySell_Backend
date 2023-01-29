package subcategory

import "context"

type SubcategoryStorage interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, SubCatID int) (int, error)
	Delete(ctx context.Context, SubCatID int) error
	GetOne(ctx context.Context, SubCatID int) (map[string]interface{}, error)
}
type ProductStorage interface {
	GetAll(ctx context.Context, SubCatID int) ([]map[string]interface{}, error)
}

type SubcategoryService struct {
	SubcatStore SubcategoryStorage
	ProdStore   ProductStorage
}

func (s *SubcategoryService) Create(ctx context.Context, m map[string]interface{}, UserID, CatID int) (int, error) {
	return s.SubcatStore.Create(ctx, m)
}

func (s *SubcategoryService) Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error) {
	return s.SubcatStore.Update(ctx, m, SubCatID)
}

func (s *SubcategoryService) Delete(ctx context.Context, UserID, CatID, SubCatID int) error {
	return s.SubcatStore.Delete(ctx, SubCatID)
}

func (s *SubcategoryService) Get(ctx context.Context, UserID, CatID, SubCatID int) (map[string]interface{}, error) {
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
