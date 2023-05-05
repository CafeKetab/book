package publishers

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
	// insert a publisher
	Insert(context.Context, *models.Publisher) error

	// get publisher detail
	GetById(ctx context.Context, id uint64) (*models.Publisher, error)

	// search + pagination (no detail)
	GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Publisher, string, error)
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func (r *repository) Insert(ctx context.Context, publisher *models.Publisher) error {
	if len(publisher.Name) == 0 || len(publisher.Title) == 0 {
		return errors.New("Insufficient information for publisher")
	}

	in := []interface{}{publisher.Name, publisher.Title, publisher.Description}
	out := []any{&publisher.Id}
	if err := r.rdbms.QueryRow(QueryInsert, in, out); err != nil {
		r.logger.Error("Error inserting publisher", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*models.Publisher, error) {
	publisher := models.Publisher{Id: id}

	out := []any{&publisher.Name, &publisher.Title, &publisher.Description}
	if err := r.rdbms.QueryRow(QueryGetDetail, []any{id}, out); err != nil {
		r.logger.Error("Error get publisher by id", zap.Error(err))
		return nil, err
	}

	return &publisher, nil
}

func (r *repository) GetAll(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Publisher, string, error) {
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

	publishers := make([]models.Publisher, limit)
	out := make([][]any, limit)

	for index := 0; index < limit; index++ {
		out[index] = []any{&publishers[index].Id, &publishers[index].Name, &publishers[index].Title}
	}

	if err := r.rdbms.Query(QueryGetAll, []any{id, search, limit}, out); err != nil {
		r.logger.Error("Error query publishers", zap.Error(err))
		return nil, "", err
	}

	var lastPublisher models.Publisher

	for index := limit - 1; index >= 0; index-- {
		if publishers[index].Id != 0 {
			lastPublisher = publishers[index]
			break
		} else {
			publishers = publishers[:index]
		}
	}

	if lastPublisher.Id == 0 {
		return publishers, "", nil
	}

	cursor := strconv.FormatUint(lastPublisher.Id, 10)

	// encrypt cursor
	encryptedCursor, err := crypto.Encrypt(cursor, r.config.CursorSecret)
	if err != nil {
		panic(err)
	}

	return publishers, encryptedCursor, nil
}
