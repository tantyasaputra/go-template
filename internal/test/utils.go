package test

import (
	"go-template/internal/log"

	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

// generally this error is not useful after logging
// make caller ask for the error using an out param
func ResetTestDB(d *gorm.DB, err ...error) {
	// Leverage Postgres to remove data from all FDI tables
	// Refer: https://stackoverflow.com/questions/2829158/truncating-all-tables-in-a-postgres-database
	result := d.Exec(`
		DO
		$do$
		BEGIN
		-- dangerous, test before you execute!
		-- RAISE NOTICE '%',  -- once confident, comment this line ...
		EXECUTE         -- ... and uncomment this one
		(SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
			FROM   pg_class
			WHERE  relkind = 'r'  -- only tables
			AND    relname != 'goose_db_version' -- skip goose migration table
			AND    relnamespace = 'public'::regnamespace
		);
		END
		$do$;
	`)

	log.Infow("DB Reset", "event", "reset failed", "error", result.Error)

	if len(err) > 0 {
		err[0] = result.Error
	}
}

// generally this error is not useful after logging
// make caller ask for the error using an out param
func FullResetTestDB(d *gorm.DB, err ...error) {
	// Leverage Postgres to completely reset database
	// Refer: https://stackoverflow.com/questions/2829158/truncating-all-tables-in-a-postgres-database
	result := d.Exec(`
		DO
		$do$
		BEGIN
		-- dangerous, test before you execute!
		-- RAISE NOTICE '%',  -- once confident, comment this line ...
		EXECUTE         -- ... and uncomment this one
		(SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
			FROM   pg_class
			WHERE  relkind = 'r'  -- only tables
			AND    relnamespace = 'public'::regnamespace
		);
		END
		$do$;
	`)

	log.Infow("DB Reset", "event", "reset failed", "error", result.Error)

	if len(err) > 0 {
		err[0] = result.Error
	}
}

// Convert string to pgtype.Numeric
func StringToNumeric(amountStr string) pgtype.Numeric {
	var amount pgtype.Numeric

	if err := amount.ScanScientific(amountStr); err != nil {
		log.Warnw("String to Numeric", "event", "string scan failed", "error", err, "service", "IKN_B2B")
	}

	return amount
}
