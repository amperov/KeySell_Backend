package subcategory

import "time"

type CreateSubcatInput struct {
	TitleRU       string `json:"title_ru,omitempty"`
	TitleEng      string `json:"title_eng,omitempty"`
	CatID         int    `json:"category_id"`
	SubtypeValue  int    `json:"subtype_value,omitempty"`
	PartialValues string `json:"partial_values,omitempty"`
	IsComposite   bool   `json:"is_composite,omitempty"`
	CreatedAt     time.Time
	SubitemID     int `json:"subitem_id"`
}

func (c *CreateSubcatInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["title_ru"] = c.TitleRU
	m["title_eng"] = c.TitleEng
	m["subitem_id"] = c.SubitemID
	m["is_composite"] = c.IsComposite
	m["partial_values"] = c.PartialValues
	m["subtype_value"] = c.SubtypeValue
	m["created_at"] = time.Now().String()
	m["category_id"] = c.CatID

	return m
}

type UpdateSubcatInput struct {
	TitleRU   string `json:"title_ru,omitempty"`
	TitleEng  string `json:"title_eng,omitempty"`
	CatID     int    `json:"category_id,omitempty"`
	SubItemID int    `json:"subitem_id,omitempty"`
}

func (c *UpdateSubcatInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if c.TitleRU != "" {
		m["title_ru"] = c.TitleRU
	}
	if c.TitleEng != "" {
		m["title_eng"] = c.TitleEng
	}
	if c.SubItemID != 0 {
		m["subitem_id"] = c.SubItemID
	}
	m["category_id"] = c.CatID
	return m
}
