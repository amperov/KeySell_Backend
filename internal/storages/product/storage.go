package product

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
)

type ProductStorage struct {
	c *pgxpool.Pool
}

func NewProductStorage(c *pgxpool.Pool) *ProductStorage {
	return &ProductStorage{c: c}
}

var (
	prodTable = "products"
)

func (p *ProductStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(prodTable).PlaceholderFormat(squirrel.Dollar).SetMap(m).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}
	row := p.c.QueryRow(ctx, query, args...)

	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (p *ProductStorage) Update(ctx context.Context, m map[string]interface{}, ProdID int) (int, error) {

	query, args, err := squirrel.Update(prodTable).PlaceholderFormat(squirrel.Dollar).Set("content_key", m["content_key"]).Where(squirrel.Eq{"id": ProdID}).ToSql()
	if err != nil {
		return 0, err
	}
	_, err = p.c.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return ProdID, nil
}

func (p *ProductStorage) Delete(ctx context.Context, ProdID int) error {
	query, args, err := squirrel.Delete(prodTable).PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"id": ProdID}).ToSql()
	if err != nil {
		return err
	}
	_, err = p.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductStorage) GetAll(ctx context.Context, CatID, SubCatID int) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	var i ProdForClient
	query, args, err := squirrel.Select("id", "content_key", "created_at").From(prodTable).
		Where(squirrel.Eq{"subcategory_id": SubCatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		i.SubcategoryID = SubCatID
		err := rows.Scan(&i.ID, &i.Content, &i.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		m = append(m, i.ToMapForSeller())
	}
	return m, nil
}

func (p *ProductStorage) GetCount(ctx context.Context, SubCatID int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE subcategory_id=$1", prodTable)

	row := p.c.QueryRow(ctx, query, SubCatID)
	err := row.Scan(&count)
	if err != nil {
		logrus.Println("GetCount: ", err)
		return 0, err
	}

	return count, nil
}
func (p *ProductStorage) GetCountForSelectTool(ctx context.Context, SubCatID int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE subcategory_id=$1", prodTable)

	row := p.c.QueryRow(ctx, query, SubCatID, false)
	err := row.Scan(&count)
	if err != nil {
		logrus.Println("GetCountForSelectTool: ", err)
		return 0, err
	}

	return count, nil
}

func (p *ProductStorage) GetForClient(ctx context.Context, SubcatID, Count int) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	var i ProdForClient

	query, args, err := squirrel.Select("content_key", "id", "created_at").PlaceholderFormat(squirrel.Dollar).From(prodTable).
		Where(squirrel.Eq{"subcategory_id": SubcatID}).Suffix(fmt.Sprintf("LIMIT %d", Count)).ToSql()
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&i.Content, &i.ID, &i.CreatedAt)
		if err != nil {
			logrus.Println(err)
			return nil, err
		}
		i.SubcategoryID = SubcatID

		m = append(m, i.ToMapForSeller())
	}

	return m, nil
}

func (p *ProductStorage) DeleteOne(ctx context.Context, ProdID int) error {
	query, args, err := squirrel.Delete(prodTable).Where(squirrel.Eq{"id": ProdID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = p.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (p *ProductStorage) PreCheck(ctx context.Context, SubItemID int) bool {
	var exists bool

	var id int
	var subcatStr string
	query, args, err := squirrel.Select("id", "title_ru").PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"subitem_id": SubItemID}).From("subcategory").ToSql()
	if err != nil {
		return false
	}

	row := p.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id, &subcatStr)

	if err != nil {
		log.Println(err)
		return false
	}

	query, args, err = squirrel.Select("id").Prefix("SELECT EXISTS(").Suffix(")").From(prodTable).
		Where(squirrel.Eq{"subcategory_id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return false
	}

	row = p.c.QueryRow(ctx, query, args...)

	err = row.Scan(&exists)
	if err != nil {
		log.Println("Scan prods err: ", err)
		return false
	}
	return exists
}
