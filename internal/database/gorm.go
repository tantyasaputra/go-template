package database

import (
	"context"
	"database/sql"
	"go-template/internal/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormHandler represents db handler for gorm
type GormHandler struct {
	db  *gorm.DB
	sql *sql.DB
}

// DataHandler is interface for Gorm Handler
type DataHandler interface {
	// GetDB Return DB or TX & Force WithContext
	GetDB(ctx context.Context) *gorm.DB
	// RunTransaction : All repo call using the context given will included in transaction, returning error will result in rollback while nil will result in commit.
	RunTransaction(ctx context.Context, fc func(ctx context.Context) error) error

	// Ping checks the db is still responsive. Useful for health checks.
	Ping(ctx context.Context) error
}

type contextKey string

const txkey = contextKey("DBTX")

// NewDataHandler for repository, automatically include migration with the same db
func NewDataHandler(dbconn string) DataHandler {
	// Initiate DB
	db, err := gorm.Open(postgres.Open(dbconn), &gorm.Config{})
	if err != nil {
		log.Fatalw("Failed to connect database", "event", "database init", "error", err)
	}

	mig, err := db.DB()
	if err != nil {
		log.Fatalw("Failed to load migration connection", "event", "migrator init", "error", err)
	}
	// Return Handler
	return &GormHandler{
		db:  db,
		sql: mig,
	}
}

// GetDB Return DB or TX & Force WithContext
func (h *GormHandler) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txkey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return h.db.WithContext(ctx)
}

func (h *GormHandler) getRawDB() *sql.DB {
	return h.sql
}

func (h *GormHandler) Ping(ctx context.Context) error {
	return h.sql.PingContext(ctx)
}

// RunTransaction : All repo call using the context given will included in transaction, returning error will result in rollback while nil will result in commit.
func (h *GormHandler) RunTransaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, txkey, tx))
	})
}
