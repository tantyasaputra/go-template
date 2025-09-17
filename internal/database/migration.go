package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go-template/internal/log"
	"io"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
)

const (
	commandUpConst      = "up"
	commandDownConst    = "down"
	commandStatusConst  = "status"
	commandResetConst   = "reset"
	commandVersionConst = "version"
	rootServer          = "/bin/go-template" // set the root server as the BU name
	schemaDir           = "/migrations"
	seedDir             = "/migrations/seeds"
)

type MigrationHandler interface {
	Up() error
	Down() error
	Status() error
	Reset() error
	Version() error
	SetVerbose(bool)
}

type GooseHandler struct {
	db  *sql.DB
	dir string
}

func NewGooseHandler(data *GormHandler) MigrationHandler {
	goose.SetLogger(log.Goose)

	return &GooseHandler{
		db:  data.getRawDB(),
		dir: getMigrationPath(),
	}
}

func (m *GooseHandler) Up() error {
	return m.runGoose(commandUpConst)
}

func (m *GooseHandler) Down() error {
	return m.runGoose(commandDownConst)
}

func (m *GooseHandler) Status() error {
	return m.runGoose(commandStatusConst)
}

func (m *GooseHandler) Reset() error {
	return m.runGoose(commandResetConst)
}

func (m *GooseHandler) Version() error {
	return m.runGoose(commandVersionConst)
}

func (m *GooseHandler) SetVerbose(addVerbosity bool) {
	goose.SetVerbose(addVerbosity)
}

func (m *GooseHandler) runGoose(command string) error {
	if err := goose.Run(command, m.db, fmt.Sprint(m.dir, schemaDir), ""); err != nil {
		return err
	}

	seed_path := fmt.Sprint(m.dir, seedDir)
	empty, err := dirIsEmpty(seed_path)

	if err != nil {
		return err
	}

	// do seeds as well
	if !empty {
		// skip if command is to reset
		if strings.EqualFold(command, commandResetConst) {
			return nil
		}

		if err := goose.RunWithOptions(command, m.db, seed_path, make([]string, 0), goose.WithNoVersioning()); err != nil {
			return err
		}
	}

	return nil
}

func getMigrationPath() string {
	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")
	for len(path) > 0 {
		if _, err := os.Stat(path + schemaDir); !os.IsNotExist(err) {
			// migration found
			return path
		}
		// if not found, go up 1 directory
		// check if it cant go up anymore, force exit
		if !strings.Contains(path, "/") {
			return ""
		}
		// look for last index of / in path
		idx := strings.LastIndex(path, "/")
		// slice the string
		path = path[0:idx]
	}
	return path
}

func dirIsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// readme is in the directory but we don't include it in the script count
	n, err := f.Readdirnames(2)

	if errors.Is(err, io.EOF) || len(n) < 2 {
		return true, nil
	}

	return false, err // Either not empty or error, suits both cases
}
