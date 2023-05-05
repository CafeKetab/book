package categories

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/CafeKetab/book/internal/models"
	"github.com/CafeKetab/book/pkg/crypto"
	"github.com/CafeKetab/book/pkg/rdbms"
	"go.uber.org/zap"
)

type Repository interface {
	// insert a category
	Insert(context.Context, *models.Category) error

	// get category detail
	GetById(ctx context.Context, id uint64) (*models.Category, error)

	// search + pagination (no detail)
	GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Category, string, error)
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func (r *repository) Insert(ctx context.Context, category *models.Category) error {
	if len(category.Name) == 0 || len(category.Title) == 0 {
		return errors.New("Insufficient information for category")
	}

	in := []interface{}{category.Name, category.Title, category.Description}
	out := []any{&category.Id}
	if err := r.rdbms.QueryRow(QueryInsert, in, out); err != nil {
		r.logger.Error("Error inserting category", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*models.Category, error) {
	category := models.Category{Id: id}

	out := []any{&category.Name, &category.Title, &category.Description}
	if err := r.rdbms.QueryRow(QueryGetDetail, []any{id}, out); err != nil {
		r.logger.Error("Error find category by id", zap.Error(err))
		return nil, err
	}

	return &category, nil
}

func (r *repository) GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Category, string, error) {
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

	if err := r.rdbms.Query(QueryGetAll, []any{id, search, limit}, out); err != nil {
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
