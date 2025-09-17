package example

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-template/internal/repository/example"
	mocks "go-template/mocks/repository/example"
)

type UnitTestSuite struct {
	suite.Suite
	exampleRepoMock *mocks.Repository
	service         Service
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

// Will run every test as initiator
func (uts *UnitTestSuite) SetupTest() {
	// to avoid overlapping mock return on our mocks, we need to re-define it every time test is run
	// Initiate mock repository
	uts.exampleRepoMock = &mocks.Repository{}
	// Initiate service using mock repository
	uts.service = NewExampleService(uts.exampleRepoMock)
}

// Test Case 1 - Get 2 data
func (uts *UnitTestSuite) TestExampleGet() {
	want := []*example.Example{
		{
			ID:   1,
			Name: "a",
		},
		{
			ID:   2,
			Name: "a",
		},
	}
	// GetRecord (ctx context.Context, ids []int) ([]*example.Example, error)
	uts.exampleRepoMock.On("GetRecords", context.TODO(), []int{}).Return([]*example.Example{
		{
			ID:   1,
			Name: "a",
		},
		{
			ID:   2,
			Name: "a",
		},
	}, nil)
	// trigger service to test
	actual, err := uts.service.ExampleGet(context.TODO())

	uts.Equal(want, actual)
	uts.Nil(err)
}

// Test Case 2 - Error occurs when hitting repo
func (uts *UnitTestSuite) TestExampleGet_ErrorFromRepo() {
	expectedError := errors.New("internal server error")
	uts.exampleRepoMock.On("GetRecords", context.TODO(), []int{}).Return(nil, errors.New("failed to get data"))
	actual, err := uts.service.ExampleGet(context.TODO())

	uts.Nil(actual)
	uts.Equal(expectedError, err)
}
