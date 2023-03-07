package subcategory

type Subcategory struct {
	ID            int
	TitleRU       string
	TitleENG      string
	SubItemID     int
	CategoryID    int
	CreateAt      string
	SubtypeValue  int
	PartialValues string
	IsComposite   bool
}

func (s *Subcategory) ToMap() map[string]interface{} {
	ModelMap := make(map[string]interface{})
	ModelMap["id"] = s.ID
	ModelMap["subtype_value"] = s.SubtypeValue
	ModelMap["partial_values"] = s.PartialValues
	ModelMap["is_composite"] = s.IsComposite
	ModelMap["title_ru"] = s.TitleRU
	ModelMap["title_eng"] = s.TitleENG
	ModelMap["subitem_id"] = s.SubItemID
	ModelMap["category_id"] = s.CategoryID
	ModelMap["created_at"] = s.CreateAt
	return ModelMap
}
