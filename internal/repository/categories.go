package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/CafeKetab/book/internal/models"
	"github.com/CafeKetab/book/pkg/crypto"
	"go.uber.org/zap"
)

const QueryGetCategoryDetail = `
	SELECT name, title, description 
	FROM categories 
	WHERE id=$1;`

// get category detail
func (r *repository) GetCategoryById(ctx context.Context, id uint64) (*models.Category, error) {
	category := models.Category{}

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
		name LIKE '%$2%'
	ORDER BY id
	FETCH NEXT $3 ROWS ONLY;`

// search + pagination (no detail)
func (r *repository) GetCategories(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Category, string, error) {
	var id uint64 = 0

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

	if err := r.rdbms.Query(QueryGetCategoryDetail, []any{id, search, limit}, out); err != nil {
		r.logger.Error("Error query categories", zap.Error(err))
		return nil, "", err
	}

	var lastCategory models.Category

	for index := limit - 1; index >= 0; index-- {
		if categories[index].Id != 0 {
			lastCategory = categories[index]
			break
		}
	}

	if lastCategory.Id != 0 {
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
