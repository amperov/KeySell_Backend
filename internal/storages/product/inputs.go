package product

type ProdForClient struct {
	ID            int
	Content       string
	Category      string
	Subcategory   string
	DateCheck     string
	SubcategoryID int
}

func (c *ProdForClient) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = c.ID
	m["content"] = c.Content
	m["category"] = c.Category
	m["subcategory"] = c.Subcategory
	m["date_check"] = c.DateCheck
	return m
}
func (c *ProdForClient) ToMapForSeller() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = c.ID
	m["content"] = c.Content
	m["subcategory_id"] = c.SubcategoryID
	return m
}
