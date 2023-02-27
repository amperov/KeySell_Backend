package subcategory

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"sort"
)

var table = "subcategory"

type SubcategoryStorage struct {
	c *pgxpool.Pool
}

func NewSubcategoryStorage(c *pgxpool.Pool) *SubcategoryStorage {
	return &SubcategoryStorage{c: c}
}

func (c *SubcategoryStorage) GetAll(ctx context.Context, CatID int) ([]map[string]interface{}, error) {
	var subcats []Subcategory
	query, args, err := squirrel.Select("id", "title_ru", "title_eng", "subitem_id").Where(squirrel.Eq{"category_id": CatID}).PlaceholderFormat(squirrel.Dollar).From(table).ToSql()
	if err != nil {
		logrus.Printf("error make query: %v", err)
		return nil, err
	}

	rows, err := c.c.Query(ctx, query, args...)
	if err != nil {
		logrus.Printf("error query: %v", err)
		return nil, err
	}
	var arrayMap []map[string]interface{}
	for rows.Next() {
		var cat Subcategory

		err = rows.Scan(&cat.ID, &cat.TitleRU, &cat.TitleENG, &cat.SubItemID)
		if err != nil {
			logrus.Printf("error scan: %v", err)
			return nil, err
		}
		cat.CategoryID = CatID
		subcats = append(subcats, cat)
	}

	sort.Slice(subcats, func(i, j int) bool {
		return subcats[i].SubItemID < subcats[j].SubItemID
	})

	for _, subcategory := range subcats {
		arrayMap = append(arrayMap, subcategory.ToMap())
	}

	return arrayMap, nil
}

func (c *SubcategoryStorage) GetCount(ctx context.Context, CatID int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE category_id=$1", table)

	row := c.c.QueryRow(ctx, query, CatID)

	err := row.Scan(&count)
	if err != nil {
		logrus.Printf("error: %v", err)
		return 0, err
	}
	return count, nil
}

func (c *SubcategoryStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(table).SetMap(m).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		logrus.Printf("error: %v", err)
		return 0, err
	}
	return id, nil
}

func (c *SubcategoryStorage) Update(ctx context.Context, m map[string]interface{}, SubCatID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": SubCatID}).
		PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").SetMap(m).ToSql()
	if err != nil {
		logrus.Printf("error: %v", err)
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		logrus.Printf("error: %v", err)
		return 0, err
	}
	return id, nil
}

func (c *SubcategoryStorage) Delete(ctx context.Context, SubCatID int) error {
	query, args, err := squirrel.Delete(table).Where(squirrel.Eq{"id": SubCatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		logrus.Printf("error: %v", err)
		return err
	}
	_, err = c.c.Exec(ctx, query, args...)
	if err != nil {
		logrus.Printf("error: %v", err)
		return err
	}
	return nil
}

func (c *SubcategoryStorage) GetOne(ctx context.Context, SubCatID int) (map[string]interface{}, error) {
	var cat Subcategory

	query, args, err := squirrel.Select("title_ru", "title_eng", "category_id", "subitem_id").Where(squirrel.Eq{"id": SubCatID}).From(table).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		logrus.Printf("error: %v", err)
		return nil, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&cat.TitleRU, &cat.TitleENG, &cat.CategoryID, &cat.SubItemID)
	if err != nil {
		logrus.Printf("error: %v", err)
		return nil, err
	}

	cat.ID = SubCatID
	return cat.ToMap(), nil
}

func (c *SubcategoryStorage) GetIDByTitle(ctx context.Context, Title string, CategoryID int) (int, error) {
	var id int
	query := "SELECT id FROM subcategory WHERE (title_ru=$1 AND category_id=$2) OR (title_eng=$1 AND category_id=$2)"

	row := c.c.QueryRow(ctx, query, Title, CategoryID)

	err := row.Scan(&id)
	if err != nil {
		logrus.Debugf("Subcat error: %v", err)
		return 0, err
	}
	return id, nil
}
