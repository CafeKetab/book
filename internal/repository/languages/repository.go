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
	// insert a language
	Insert(context.Context, *models.Language) error

	// get language detail
	GetById(ctx context.Context, id uint64) (*models.Language, error)

	// search + pagination (no detail)
	GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Language, string, error)
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func (r *repository) Insert(ctx context.Context, category *models.Language) error {
	if len(category.Name) == 0 {
		return errors.New("Insufficient information for language")
	}

	in := []interface{}{category.Name}
	out := []any{&category.Id}
	if err := r.rdbms.QueryRow(QueryInsert, in, out); err != nil {
		r.logger.Error("Error inserting language", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*models.Language, error) {
	language := models.Language{Id: id}

	out := []any{&language.Name}
	if err := r.rdbms.QueryRow(QueryGetDetail, []any{id}, out); err != nil {
		r.logger.Error("Error get language by id", zap.Error(err))
		return nil, err
	}

	return &language, nil
}

func (r *repository) GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Language, string, error) {
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

	languages := make([]models.Language, limit)
	out := make([][]any, limit)

	for index := 0; index < limit; index++ {
		out[index] = []any{&languages[index].Id, &languages[index].Name}
	}

	if err := r.rdbms.Query(QueryGetAll, []any{id, search, limit}, out); err != nil {
		r.logger.Error("Error query languages", zap.Error(err))
		return nil, "", err
	}

	var lastLanguage models.Language

	for index := limit - 1; index >= 0; index-- {
		if languages[index].Id != 0 {
			lastLanguage = languages[index]
			break
		} else {
			languages = languages[:index]
		}
	}

	if lastLanguage.Id == 0 {
		return languages, "", nil
	}

	cursor := strconv.FormatUint(lastLanguage.Id, 10)

	// encrypt cursor
	encryptedCursor, err := crypto.Encrypt(cursor, r.config.CursorSecret)
	if err != nil {
		panic(err)
	}

	return languages, encryptedCursor, nil
}
