package sample

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type (

	// Sample represents rows in table Samples
	Sample struct {
		ID     int
		Data   JSONB
		Detail Detail `gorm:"-:all"` // JSON Field Ignore on gorm call
	}

	// Detail represents json in Sample
	Detail struct {
		NetPrice    int64  `gorm:"-:all" json:"netPrice"`    // JSON Field Ignore on gorm call
		SellPrice   int64  `gorm:"-:all" json:"sellPrice"`   // JSON Field Ignore on gorm call
		Fee         int64  `gorm:"-:all" json:"fee"`         // JSON Field Ignore on gorm call
		ProductType string `gorm:"-:all" json:"productType"` // JSON Field Ignore on gorm call
	}
)

const tableJSON = "public.go_template_json"

// TableName (gorm customization): Specify table name for gorm, gorm by default will look for pluralized name from struct
func (Sample) TableName() string {
	return tableJSON
}

// Repository is the interface for Example Repository
type Repository interface {
	GetRecords(ctx context.Context, ids []int) ([]*Sample, error)
	AddRecords(ctx context.Context, payload []*Sample) error
}

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB map[string]interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
