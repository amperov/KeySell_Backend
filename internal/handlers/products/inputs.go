package products

import "time"

type CreateProductInput struct {
	Content  string `json:"content,omitempty"`
	SubCatID int    `json:"subcategory_id,omitempty"`
}

func (c *CreateProductInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["content_key"] = c.Content
	m["created_at"] = time.Now().String()
	m["subcategory_id"] = c.SubCatID
	return m
}

type UpdateProductInput struct {
	Content  string `json:"content,omitempty"`
	SubCatID int    `json:"subcategory_id,omitempty"`
}

func (c *UpdateProductInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	if c.Content != "" {
		m["content_key"] = c.Content
		return m
	}
	return nil
}

type CreateProductsInput struct {
	Products []CreateProductInput `json:"products,omitempty"`
}
