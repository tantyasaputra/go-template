package main

import (
	"flag"
	"fmt"
	"go-template/internal/config"
	"go-template/internal/database"
	"go-template/internal/log"
	"os"
	"strings"
)

var (
	verbose    bool
	migrateCmd = flag.NewFlagSet("migrate", flag.ExitOnError)
	mapActions = map[string]struct{}{
		"up":      {},
		"down":    {},
		"reset":   {},
		"status":  {},
		"version": {},
	}
)

func setupFlags() {
	for _, fs := range []*flag.FlagSet{migrateCmd} {
		fs.BoolVar(&verbose, "verbose", false, "Enable Detailed Migration Logging")
		fs.BoolVar(&verbose, "v", false, "Enable Detailed Migration Logging")
	}

	flag.Usage = func() {
		fmt.Printf("Usage: %s migrate [OPTIONS] [ACTION] \n", os.Args[0][strings.LastIndex(os.Args[0], "/")+1:])
		fmt.Println("Available Options: ")
		migrateCmd.PrintDefaults()
		fmt.Println("Available Actions: ")
		fmt.Printf("  up      \t\tMigrate the DB to the most recent version available\n")
		fmt.Printf("  down    \t\tRoll back the version by 1\n")
		fmt.Printf("  reset   \t\tRoll back all migrations\n")
		fmt.Printf("  status  \t\tDump the migration status for the current DB\n")
		fmt.Printf("  version \t\tPrint the current version of the database\n")
		fmt.Printf("\nFind this help by passing -h or --help\n")
	}
}

func handleArguments() {
	setupFlags()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := migrateCmd.Parse(os.Args[2:]); err != nil {
				log.Fatalw("Command Parser", "event", "migrator", "error", err)
			}

			if migrateCmd.NArg() == 0 {
				fmt.Printf("[Error] Missing Migrate Action\n\n")
				flag.Usage()
				os.Exit(0)
			}

			if _, ok := mapActions[migrateCmd.Arg(0)]; !ok {
				fmt.Printf("[Error] Unknown Migrate Action %s\n\n", migrateCmd.Arg(0))
				flag.Usage()
				os.Exit(0)
			}

			runGooseAction(migrateCmd.Arg(0), verbose)
			os.Exit(0)

		default:
			flag.Usage()
			os.Exit(0)
		}
	}
}

func runGooseAction(action string, verboseEnable bool) {
	dbHandler := getMigratonHandler(config.GetEnv().DB)
	dbHandler.SetVerbose(verboseEnable)

	err := detectAndRunAction(dbHandler, action)

	if err != nil {
		log.Errorw("Database Migration Error", "event", "migrator", "error", err)
	}

	if strings.EqualFold(config.GetEnv().DEVMODE, "true") {
		// migrate test db in dev mode

		fmt.Printf("\t***Dev Mode Detected; Applying Actions to Test DB***\n")

		dbHandler := getMigratonHandler(config.GetEnv().DBT)
		dbHandler.SetVerbose(verboseEnable)

		err := detectAndRunAction(dbHandler, action)

		if err != nil {
			log.Errorw("Database Migration Error", "event", "migrator", "error", err)
		}
	}
}

func getMigratonHandler(dbconn string) database.MigrationHandler {
	dataHandler := database.NewDataHandler(dbconn)
	dbHandler := database.NewGooseHandler(dataHandler.(*database.GormHandler))

	return dbHandler
}

func detectAndRunAction(db database.MigrationHandler, action string) error {
	var err error
	switch action {
	case "up":
		return db.Up()
	case "down":
		return db.Down()
	case "reset":
		return db.Reset()
	case "status":
		return db.Status()
	case "version":
		return db.Version()
	default:
		err = fmt.Errorf("unknown migration action %s", action)
	}

	return err
}
