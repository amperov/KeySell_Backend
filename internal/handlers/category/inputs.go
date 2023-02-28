package category

import "time"

type CreateCategoryInput struct {
	TitleRu     string `json:"title_ru,omitempty"`
	TitleEng    string `json:"title_eng,omitempty"`
	Description string `json:"description,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (c *CreateCategoryInput) ToMap() map[string]interface{} {
	var cat = make(map[string]interface{})
	cat["title_ru"] = c.TitleRu
	cat["title_eng"] = c.TitleEng
	cat["description"] = c.Description
	cat["created_at"] = time.Now().String()
	cat["user_id"] = c.UserID
	cat["message_client"] = c.Message
	return cat
}

type UpdateCategoryInput struct {
	TitleRu     string `json:"title_ru,omitempty"`
	TitleEng    string `json:"title_eng,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *UpdateCategoryInput) ToMap() map[string]interface{} {
	var cat = make(map[string]interface{})

	if c.TitleRu != "" {
		cat["title_ru"] = c.TitleRu
	}
	if c.TitleEng != "" {
		cat["title_eng"] = c.TitleEng
	}
	if c.Description != "" {
		cat["description"] = c.Description
	}
	if c.Message != "" {
		cat["message_client"] = c.Message
	}
	return cat
}
