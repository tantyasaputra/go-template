package postgres

import (
	"context"
	"testing"

	"go-template/internal/config"
	"go-template/internal/database"
	"go-template/internal/log"
	"go-template/internal/repository/example"
	"go-template/internal/test"

	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	repo example.Repository
	dh   database.DataHandler
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

// will run once per test run
func (uts *UnitTestSuite) SetupSuite() {
	// because we don't have any mocking involved, we can just run the setup once and not per test
	// initiate env
	if err := config.SetDevelopmentEnv(); err != nil {
		log.Fatal(err)
	}

	// Create Gorm Connection to Testing DB
	uts.dh = database.NewDataHandler(config.GetEnv().DBT)
	uts.repo = NewExampleRepository(uts.dh)
}

// Will run every test as initiator
func (uts *UnitTestSuite) SetupTest() {
	// when doing mocking, do here.
	// to avoid overlapping mock return on our mocks, we need to re-define it every time test is run
}

func (uts *UnitTestSuite) TearDownSuite() {
	test.ResetTestDB(uts.dh.GetDB(context.TODO()))
}

// Test Case 1 - AddRecords Success
func (uts *UnitTestSuite) TestAddRecords() {
	ctx := context.TODO()
	err := uts.repo.AddRecords(ctx, []*example.Example{
		{
			Name: "User 1",
		},
		{
			Name: "User 2",
		},
	})

	uts.Nil(err)
}

// Test Case 2 - AddRecords Error
func (uts *UnitTestSuite) TestAddRecords_Error() {
	err := uts.repo.AddRecords(context.TODO(), []*example.Example{})
	uts.NotNil(err)
}
