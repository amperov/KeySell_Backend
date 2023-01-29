package category

import "time"

type CreateCategoryInput struct {
	TitleRu     string `json:"title_ru,omitempty"`
	TitleEng    string `json:"title_eng,omitempty"`
	Description string `json:"description,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
}

func (c *CreateCategoryInput) ToMap() map[string]interface{} {
	var cat = make(map[string]interface{})
	cat["title_ru"] = c.TitleRu
	cat["title_eng"] = c.TitleEng
	cat["description"] = c.Description
	cat["created_at"] = time.Now().String()
	cat["user_id"] = c.UserID
	return cat
}

type UpdateCategoryInput struct {
	TitleRu     string `json:"title_ru,omitempty"`
	TitleEng    string `json:"title_eng,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *UpdateCategoryInput) ToMap() map[string]interface{} {
	var cat = make(map[string]interface{})

	if cat["title_ru"] != "" {
		cat["title_ru"] = c.TitleRu
	}
	if cat["title_eng"] != "" {
		cat["title_eng"] = c.TitleEng
	}
	if cat["description"] != "" {
		cat["description"] = c.Description
	}

	return cat
}
