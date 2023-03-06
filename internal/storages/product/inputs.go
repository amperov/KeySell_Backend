package product

type ProdForClient struct {
	ID            int
	Content       string
	Category      string
	Subcategory   string
	DateCheck     string
	SubcategoryID int
	CreatedAt     string
}

func (c *ProdForClient) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = c.ID
	m["content_key"] = c.Content
	m["category_name"] = c.Category
	m["subcategory_name"] = c.Subcategory
	m["date_check"] = c.DateCheck

	return m
}
func (c *ProdForClient) ToMapForSeller() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = c.ID
	m["content_key"] = c.Content
	m["subcategory_id"] = c.SubcategoryID
	m["created_at"] = c.CreatedAt
	return m
}
