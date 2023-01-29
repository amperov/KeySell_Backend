package subcategory

import "time"

type CreateSubcatInput struct {
	Title     string `json:"title,omitempty"`
	CatID     int    `json:"category_id"`
	CreatedAt time.Time
	SubitemID int `json:"subitem_id"`
}

func (c *CreateSubcatInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["title"] = c.Title
	m["subitem_id"] = c.SubitemID
	m["created_at"] = time.Now().String()
	m["category_id"] = c.CatID
	return m
}

type UpdateSubcatInput struct {
	Title string `json:"title,omitempty"`
	CatID int    `json:"category_id"`
}

func (c *UpdateSubcatInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["title"] = c.Title
	m["category_id"] = c.CatID
	return m
}
