package postgres

import (
	"context"
	"encoding/json"
	"go-template/internal/database"
	"go-template/internal/log"
	"go-template/internal/repository/sample"
)

type handler struct {
	db database.DataHandler
}

// NewSampleRepository return interface for example table repo
func NewSampleRepository(db database.DataHandler) sample.Repository {
	return &handler{
		db: db,
	}
}

// Sample GetRecords returns data based on ids, empty ids will return all
func (h *handler) GetRecords(ctx context.Context, ids []int) ([]*sample.Sample, error) {
	var err error
	var result []*sample.Sample
	tx := h.db.GetDB(ctx)
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
	for _, v := range result {
		js, _ := v.Data.Value()
		if err := json.Unmarshal(js.([]byte), &v.Detail); err != nil {
			log.Info(err)
		}
	}

	return result, nil
}

// Sample AddRecords return error when inserting failed
func (h *handler) AddRecords(ctx context.Context, payload []*sample.Sample) error {
	// Convert Detail into Json
	for _, v := range payload {
		js, _ := json.Marshal(v.Detail)
		if err := v.Data.Scan(js); err != nil {
			log.Info(err)
		}
	}
	return h.db.GetDB(ctx).Create(payload).Error
}
