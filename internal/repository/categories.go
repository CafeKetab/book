package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/CafeKetab/book/internal/models"
	"github.com/CafeKetab/book/pkg/crypto"
	"go.uber.org/zap"
)

const QueryInsertCategory = `
	INSERT INTO 
		categories(name, title, description) 
		VALUES($1, $2, $3) 
	RETURNING id;`

// insert a category
func (r *repository) InsertCategory(ctx context.Context, category *models.Category) error {
	if len(category.Name) == 0 || len(category.Title) == 0 {
		return errors.New("Insufficient information for category")
	}

	in := []interface{}{category.Name, category.Title, category.Description}
	if err := r.rdbms.QueryRow(QueryInsertCategory, in, []any{&category.Id}); err != nil {
		r.logger.Error("Error creating category", zap.Error(err))
		return err
	}

	return nil
}

const QueryGetCategoryDetail = `
	SELECT name, title, description 
	FROM categories 
	WHERE id=$1;`

// get category detail
func (r *repository) GetCategoryById(ctx context.Context, id uint64) (*models.Category, error) {
	category := models.Category{Id: id}

	in := []any{id}
	out := []any{&category.Name, &category.Title, &category.Description}
	if err := r.rdbms.QueryRow(QueryGetCategoryDetail, in, out); err != nil {
		r.logger.Error("Error find category by id", zap.Error(err))
		return nil, err
	}

	return &category, nil
}

const QueryGetCategories = `
	SELECT id, name, title 
	FROM categories 
	WHERE 
		id > $1 AND
		name LIKE '%' || $2 || '%'
	ORDER BY id
	FETCH NEXT $3 ROWS ONLY;`

// search + pagination (no detail)
func (r *repository) GetCategories(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Category, string, error) {
	var id uint64 = 0

	if limit < r.config.Limit.Min {
		limit = r.config.Limit.Min
	} else if limit > r.config.Limit.Max {
		limit = r.config.Limit.Max
	}

	// decrypt cursor
	if len(encryptedCursor) != 0 {
		cursor, err := crypto.Decrypt(encryptedCursor, r.config.CursorSecret)
		if err != nil {
			panic(err)
		}

		splits := strings.Split(cursor, ",")
		if len(splits) != 1 {
			panic("err")
		}

		id, err = strconv.ParseUint(splits[0], 10, 64)
		if err != nil {
			panic(err)
		}
	}

	categories := make([]models.Category, limit)
	out := make([][]any, limit)

	for index := 0; index < limit; index++ {
		out[index] = []any{&categories[index].Id, &categories[index].Name, &categories[index].Title}
	}

	if err := r.rdbms.Query(QueryGetCategories, []any{id, search, limit}, out); err != nil {
		r.logger.Error("Error query categories", zap.Error(err))
		return nil, "", err
	}

	var lastCategory models.Category

	for index := limit - 1; index >= 0; index-- {
		if categories[index].Id != 0 {
			lastCategory = categories[index]
			break
		} else {
			categories = categories[:index]
		}
	}

	if lastCategory.Id == 0 {
		return categories, "", nil
	}

	cursor := strconv.FormatUint(lastCategory.Id, 10)

	// encrypt cursor
	encryptedCursor, err := crypto.Encrypt(cursor, r.config.CursorSecret)
	if err != nil {
		panic(err)
	}

	return categories, encryptedCursor, nil
}
