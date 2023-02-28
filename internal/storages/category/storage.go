package category

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var table = "category"

type CategoryStorage struct {
	c *pgxpool.Pool
}

func NewCategoryStorage(c *pgxpool.Pool) *CategoryStorage {
	return &CategoryStorage{c: c}
}

func (c *CategoryStorage) Belong(ctx context.Context, UserID, CategoryID int) bool {
	var exist bool

	query := "SELECT EXISTS(title) FROM category WHERE user_id=$1 AND id=$2"

	row := c.c.QueryRow(ctx, query, UserID, CategoryID)
	err := row.Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}

func (c *CategoryStorage) GetIDByTitle(ctx context.Context, CategoryName string) (int, string, error) {
	var id int
	var Msg string
	query := "SELECT id, message_client FROM category WHERE (title_ru=$1) OR (title_eng=$1)"
	row := c.c.QueryRow(ctx, query, CategoryName)

	err := row.Scan(&id, &Msg)
	if err != nil {
		logrus.Debugf("Cat error: %v", err)
		return 0, "", err
	}
	return id, Msg, nil
}

func (c *CategoryStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(table).SetMap(m).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").ToSql()
	if err != nil {
		logrus.Println(err)
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		logrus.Println(err)
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) Update(ctx context.Context, m map[string]interface{}, CatID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": CatID}).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").SetMap(m).ToSql()
	if err != nil {
		logrus.Println(err)
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		logrus.Println(err)
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) Delete(ctx context.Context, CatID int) error {
	query, args, err := squirrel.Delete(table).Where(squirrel.Eq{"id": CatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		logrus.Println(err)
		return err
	}
	_, err = c.c.Exec(ctx, query, args...)
	if err != nil {
		logrus.Println(err)
		return err
	}

	return nil
}

func (c *CategoryStorage) GetOne(ctx context.Context, CatID int) (map[string]interface{}, error) {
	var cat Category

	query, args, err := squirrel.Select("title_ru", "title_eng", "description", "message_client", "item_id").
		Where(squirrel.Eq{"id": CatID}).From(table).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&cat.TitleRu, &cat.TitleEng, &cat.Description, &cat.Message, &cat.ItemID)
	if err != nil {
		return nil, err
	}

	cat.ID = CatID
	return cat.ToMap(), nil
}

func (c *CategoryStorage) GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error) {

	query, args, err := squirrel.Select("id", "title_ru", "title_eng", "description", "message_client", "item_id").
		Where(squirrel.Eq{"user_id": UserID}).
		PlaceholderFormat(squirrel.Dollar).From(table).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.c.Query(ctx, query, args...)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	var arrayMap []map[string]interface{}

	for rows.Next() {
		var cat Category

		err = rows.Scan(&cat.ID, &cat.TitleRu, &cat.TitleEng, &cat.Description, &cat.Message, &cat.ItemID)
		if err != nil {
			logrus.Print(err)
			return nil, err
		}
		cat.UserID = UserID
		logrus.Printf("CatID: %v", cat.ID)
		arrayMap = append(arrayMap, cat.ToMap())

	}

	return arrayMap, nil
}
