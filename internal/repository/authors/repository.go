package authors

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
	// insert an author
	Insert(context.Context, *models.Author) error

	// get author detail
	GetById(ctx context.Context, id uint64) (*models.Author, error)

	// search + pagination (no detail)
	GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Author, string, error)
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func New(logger *zap.Logger, config *Config, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: logger, config: config, rdbms: rdbms}
	return r
}

func (r *repository) Insert(ctx context.Context, author *models.Author) error {
	if len(author.FullName) == 0 {
		return errors.New("Insufficient information for author")
	}

	in := []interface{}{author.FullName}
	out := []any{&author.Id}
	if err := r.rdbms.QueryRow(QueryInsert, in, out); err != nil {
		r.logger.Error("Error inserting author", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*models.Author, error) {
	author := models.Author{Id: id}

	out := []any{&author.FullName}
	if err := r.rdbms.QueryRow(QueryGetDetail, []any{id}, out); err != nil {
		r.logger.Error("Error get author by id", zap.Error(err))
		return nil, err
	}

	return &author, nil
}

func (r *repository) GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Author, string, error) {
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

	authors := make([]models.Author, limit)
	out := make([][]any, limit)

	for index := 0; index < limit; index++ {
		out[index] = []any{&authors[index].Id, &authors[index].FullName}
	}

	if err := r.rdbms.Query(QueryGetAll, []any{id, search, limit}, out); err != nil {
		r.logger.Error("Error query authors", zap.Error(err))
		return nil, "", err
	}

	if len(authors) == 0 {
		return authors, "", nil
	}

	lastAuthor := authors[len(authors)]
	cursor := strconv.FormatUint(lastAuthor.Id, 10)

	// encrypt cursor
	encryptedCursor, err := crypto.Encrypt(cursor, r.config.CursorSecret)
	if err != nil {
		panic(err)
	}

	return authors, encryptedCursor, nil
}
