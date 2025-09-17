package database

import (
	"context"
	"go-template/internal/config"
	"go-template/internal/log"
	"go-template/internal/test"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseMigrationTestSuite struct {
	suite.Suite
	ctx         context.Context
	dataHandler DataHandler
}

func TestDatabaseMigrationTestSuitee(t *testing.T) {
	suite.Run(t, new(DatabaseMigrationTestSuite))
}

func (suite *DatabaseMigrationTestSuite) SetupTest() {
	if err := config.SetDevelopmentEnv(); err != nil {
		log.Fatal(err)
	}

	suite.ctx = context.Background()
	x := config.GetEnv().DBT
	suite.dataHandler = NewDataHandler(x)
}

func (suite *DatabaseMigrationTestSuite) SetupSubTest() {
	if err := NewGooseHandler(suite.dataHandler.(*GormHandler)).Up(); err != nil {
		log.Fatal(err)
	}
}

func (suite *DatabaseMigrationTestSuite) TearDownSubTest() {
	test.ResetTestDB(suite.dataHandler.GetDB(suite.ctx))
}

func (suite *DatabaseMigrationTestSuite) TestVersion() {
	migrationHandler := NewGooseHandler(suite.dataHandler.(*GormHandler))

	err := migrationHandler.Version()
	suite.Nil(err)
}

func (suite *DatabaseMigrationTestSuite) TestDown() {
	migrationHandler := NewGooseHandler(suite.dataHandler.(*GormHandler))

	err := migrationHandler.Down()
	suite.Nil(err)
}

func (suite *DatabaseMigrationTestSuite) TestStatus() {
	migrationHandler := NewGooseHandler(suite.dataHandler.(*GormHandler))

	err := migrationHandler.Status()
	suite.Nil(err)
}

func (suite *DatabaseMigrationTestSuite) TestReset() {
	migrationHandler := NewGooseHandler(suite.dataHandler.(*GormHandler))

	err := migrationHandler.Reset()
	suite.Nil(err)
}

func (suite *DatabaseMigrationTestSuite) TestUp() {
	migrationHandler := NewGooseHandler(suite.dataHandler.(*GormHandler))

	err := migrationHandler.Up()
	suite.Nil(err)
}
