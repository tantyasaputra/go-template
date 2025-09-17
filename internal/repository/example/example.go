package example

import "context"

// Example represents rows in table Examples
type Example struct {
	ID   int
	Name string
}

const tableExample = "public.go_template_examples"

// TableName (gorm customization): Specify table name for gorm, gorm by default will look for pluralized name from struct
func (Example) TableName() string {
	return tableExample
}

// Repository is the interface for Example Repository
type Repository interface {
	GetRecords(ctx context.Context, ids []int) ([]*Example, error)
	AddRecords(ctx context.Context, payload []*Example) error
}
