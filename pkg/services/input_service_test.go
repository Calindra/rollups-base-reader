package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/cartesi/rollups-graphql/pkg/commons"
	cModel "github.com/cartesi/rollups-graphql/pkg/convenience/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mocks

type MockInputRepository struct {
	mock.Mock
}

// CountMap implements repository.InputRepositoryInterface.
func (m *MockInputRepository) CountMap(ctx context.Context) (map[int64]uint64, error) {
	panic("unimplemented")
}

// Count implements repository.InputRepositoryInterface.
func (m *MockInputRepository) Count(ctx context.Context, filter []*cModel.ConvenienceFilter) (uint64, error) {
	panic("unimplemented")
}

// Create implements repository.InputRepositoryInterface.
func (m *MockInputRepository) Create(ctx context.Context, input model.Input) error {
	args := m.Called(input)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// FindAll implements repository.InputRepositoryInterface.
func (m *MockInputRepository) FindAll(ctx context.Context, first *int, last *int, after *string, before *string, filter []*cModel.ConvenienceFilter) (*commons.PageResult[model.Input], error) {
	panic("unimplemented")
}

// FindOne implements repository.InputRepositoryInterface.
func (m *MockInputRepository) FindOne(ctx context.Context, applicationId int64, index uint64) (*model.Input, error) {
	panic("unimplemented")
}

type MockEpochRepository struct {
	mock.Mock
}

// Create implements repository.EpochRepositoryInterface.
func (m *MockEpochRepository) Create(ctx context.Context, epoch *model.Epoch) error {
	panic("unimplemented")
}

// FindOne implements repository.EpochRepositoryInterface.
func (m *MockEpochRepository) FindOne(ctx context.Context, index uint64) (*model.Epoch, error) {
	panic("unimplemented")
}

// GetLatestOpenEpochByAppID implements repository.EpochRepositoryInterface.
func (m *MockEpochRepository) GetLatestOpenEpochByAppID(ctx context.Context, appID int64) (*model.Epoch, error) {
	args := m.Called(appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Epoch), args.Error(1)
}

type MockAppRepository struct {
	mock.Mock
}

// FindOneByID implements repository.AppRepositoryInterface.
func (m *MockAppRepository) FindOneByID(ctx context.Context, id int64) (*model.Application, error) {
	panic("unimplemented")
}

// FindAllByDA implements repository.AppRepositoryInterface.
func (m *MockAppRepository) FindAllByDA(ctx context.Context, da model.DataAvailabilitySelector) ([]model.Application, error) {
	args := m.Called(da)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Application), args.Error(1)
}

// FindOneByContract implements repository.AppRepositoryInterface.
func (m *MockAppRepository) FindOneByContract(ctx context.Context, address common.Address) (*model.Application, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Application), args.Error(1)
}

// List implements repository.AppRepositoryInterface.
func (m *MockAppRepository) List(ctx context.Context) ([]model.Application, error) {
	panic("unimplemented")
}

// UpdateDA implements repository.AppRepositoryInterface.
func (m *MockAppRepository) UpdateDA(ctx context.Context, applicationId int64, da model.DataAvailabilitySelector) error {
	panic("unimplemented")
}

// Tests

const timeout = 5 * time.Minute

type InputServiceTestSuite struct {
	suite.Suite
	mockInputRepo *MockInputRepository
	mockEpochRepo *MockEpochRepository
	mockAppRepo   *MockAppRepository
	service       *InputService
	ctx           context.Context
	ctxCancel     context.CancelFunc
}

func TestInputService(t *testing.T) {
	suite.Run(t, new(InputServiceTestSuite))
}

func (is *InputServiceTestSuite) SetupTest() {
	is.ctx, is.ctxCancel = context.WithTimeout(context.Background(), timeout)

	is.mockInputRepo = new(MockInputRepository)
	is.mockEpochRepo = new(MockEpochRepository)
	is.mockAppRepo = new(MockAppRepository)

	is.service = NewInputService(is.mockInputRepo, is.mockEpochRepo, is.mockAppRepo)
}

func (is *InputServiceTestSuite) TearDownTest() {
	is.ctxCancel()
}

func (is *InputServiceTestSuite) TestCreateInputID() {
	ctx, ctxCancel := context.WithCancel(is.ctx)
	defer ctxCancel()

	appID := int64(1)
	input := model.Input{RawData: []byte("test-payload"), EpochApplicationID: appID, EpochIndex: 999}

	mockEpoch := &model.Epoch{Index: 10, Status: model.EpochStatus_Open, ApplicationID: appID}

	updatedInput := model.Input{
		RawData:            input.RawData,
		EpochApplicationID: appID,
		EpochIndex:         mockEpoch.Index,
	}

	is.mockEpochRepo.On("GetLatestOpenEpochByAppID", appID).Return(mockEpoch, nil)
	is.mockInputRepo.On("Create", updatedInput).Return(nil)

	err := is.service.CreateInput(ctx, input)

	is.NoError(err)
	is.mockEpochRepo.AssertCalled(is.T(), "GetLatestOpenEpochByAppID", appID)
	is.mockInputRepo.AssertCalled(is.T(), "Create", updatedInput)
}

func (is *InputServiceTestSuite) TestCreateInputWithAddress() {
	ctx, ctxCancel := context.WithCancel(is.ctx)
	defer ctxCancel()

	appContract := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	input := model.Input{RawData: []byte("test-payload"), EpochApplicationID: 0, EpochIndex: 0}

	mockApp := &model.Application{ID: 13}
	is.mockAppRepo.On("FindOneByContract", appContract).Return(mockApp, nil)
	mockEpoch := &model.Epoch{Index: 10, ApplicationID: mockApp.ID, Status: model.EpochStatus_Open}

	updatedInput := model.Input{
		RawData:            input.RawData,
		EpochApplicationID: mockApp.ID,
		EpochIndex:         mockEpoch.Index,
	}

	is.mockEpochRepo.On("GetLatestOpenEpochByAppID", mockApp.ID).Return(mockEpoch, nil)
	is.mockInputRepo.On("Create", updatedInput).Return(nil)

	err := is.service.CreateInputWithAddress(ctx, appContract, input)

	is.NoError(err)
	is.mockAppRepo.AssertCalled(is.T(), "FindOneByContract", appContract)
	is.mockEpochRepo.AssertCalled(is.T(), "GetLatestOpenEpochByAppID", mockApp.ID)
	is.mockInputRepo.AssertCalled(is.T(), "Create", updatedInput)
}

func (is *InputServiceTestSuite) TestCreateInputWithAddress_AppNotFound() {
	ctx, ctxCancel := context.WithCancel(is.ctx)
	defer ctxCancel()

	appContract := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	input := model.Input{RawData: []byte("test-payload")}

	is.mockAppRepo.On("FindOneByContract", appContract).Return((*model.Application)(nil), errors.New("app not found"))

	err := is.service.CreateInputWithAddress(ctx, appContract, input)

	is.Error(err)
	is.mockAppRepo.AssertCalled(is.T(), "FindOneByContract", appContract)
	is.mockEpochRepo.AssertNotCalled(is.T(), "GetLatestOpenEpochByAppID", mock.Anything, mock.Anything)
	is.mockInputRepo.AssertNotCalled(is.T(), "Create", mock.Anything, mock.Anything)
}
