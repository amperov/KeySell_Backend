package subcategory

type Subcategory struct {
	ID         int
	TitleRU    string
	TitleENG   string
	SubItemID  int
	CategoryID int
	CreateAt   string
}

func (s *Subcategory) ToMap() map[string]interface{} {
	ModelMap := make(map[string]interface{})
	ModelMap["id"] = s.ID

	ModelMap["title_ru"] = s.TitleRU
	ModelMap["title_eng"] = s.TitleENG
	ModelMap["subitem_id"] = s.SubItemID
	ModelMap["category_id"] = s.CategoryID
	ModelMap["created_at"] = s.CreateAt
	return ModelMap
}
