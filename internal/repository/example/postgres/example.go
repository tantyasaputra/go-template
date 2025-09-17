package postgres

import (
	"context"
	"go-template/internal/database"
	"go-template/internal/repository/example"
)

type handler struct {
	db database.DataHandler
}

// NewExampleRepository return interface for example table repo
func NewExampleRepository(db database.DataHandler) example.Repository {
	return &handler{
		db: db,
	}
}

// Example GetRecords returns data based on ids, empty ids will return all
func (h *handler) GetRecords(ctx context.Context, ids []int) ([]*example.Example, error) {
	var err error
	var result []*example.Example
	tx := h.db.GetDB(ctx).Where("deleted_at is null")
	if len(ids) > 0 {
		// get based on ids
		tx = tx.Find(&result, ids)
	} else {
		// get all
		tx = tx.Find(&result)
	}
	// check error
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Example AddRecords return error when inserting failed
func (h *handler) AddRecords(ctx context.Context, payload []*example.Example) error {
	return h.db.GetDB(ctx).Create(payload).Error
}
