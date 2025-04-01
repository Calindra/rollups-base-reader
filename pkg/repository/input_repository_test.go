package repository

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/calindra/rollups-base-reader/pkg/model"
	cModel "github.com/cartesi/rollups-graphql/pkg/convenience/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type InputRepositorySuite struct {
	suite.Suite
	inputRepository *InputRepository
	ctx             context.Context
	ctxCancel       context.CancelFunc
	image           *postgres.PostgresContainer
	schemaDir       string
}

func (s *InputRepositorySuite) SetupSuite() {
	// Fetch schema
	tmpDir, err := os.MkdirTemp("", "schema")
	s.NoError(err)
	s.schemaDir = filepath.Join(tmpDir, "schema.sql")
	schemaFile, err := os.Create(s.schemaDir)
	s.NoError(err)
	defer schemaFile.Close()

	resp, err := http.Get(commons.Schema)
	s.NoError(err)
	defer resp.Body.Close()

	_, err = io.Copy(schemaFile, resp.Body)
	s.NoError(err)
}

func (s *InputRepositorySuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), commons.DefaultTimeout)

	// Database
	container, err := postgres.Run(s.ctx, commons.DbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaDir),
		postgres.WithDatabase(commons.DbName),
		postgres.WithUsername(commons.DbUser),
		postgres.WithPassword(commons.DbPassword),
		testcontainers.WithLogConsumers(&commons.StdoutLogConsumer{}),
	)
	s.NoError(err)
	extraArg := "sslmode=disable"
	connectionStr, err := container.ConnectionString(s.ctx, extraArg)
	s.NoError(err)
	s.image = container
	err = container.Start(s.ctx)
	s.NoError(err)

	db, err := sqlx.ConnectContext(s.ctx, "postgres", connectionStr)
	s.NoError(err)

	s.inputRepository = NewInputRepository(db)
}

func (s *InputRepositorySuite) TearDownTest() {
	err := s.image.Stop(s.ctx, nil)
	s.NoError(err)
	s.inputRepository.Db.Close()
	s.ctxCancel()
}

func (s *InputRepositorySuite) TestInputRepository() {
	input := model.Input{
		EpochApplicationID: 1,         // existing app
		EpochIndex:         commons.OpenEpoch, // add to actual epoch
		Index:              171,       // unique index
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             model.InputCompletionStatus_Accepted,
	}
	err := s.inputRepository.Create(s.ctx, input)
	s.NoError(err)

	inputDb, err := s.inputRepository.FindOne(s.ctx, input.EpochApplicationID, input.Index)
	s.NoError(err)
	s.Equal(input.EpochApplicationID, inputDb.EpochApplicationID)
}

func (s *InputRepositorySuite) TestInputWrongIndex() {
	input := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              1,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             model.InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestInputWrongEpoch() {
	input := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         999, // non-existent epoch
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             model.InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestInputWrongApplication() {
	input := model.Input{
		EpochApplicationID: 999, // non-existent application
		EpochIndex:         commons.OpenEpoch,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             model.InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestQueryInputWrongApplicationIndex() {
	var (
		applicationId int64  = 999 // non-existent application
		index         uint64 = 171 // non-existent index
	)
	_, err := s.inputRepository.FindOne(s.ctx, applicationId, index)
	s.Error(err)
	s.ErrorIs(err, sql.ErrNoRows)
}

func (s *InputRepositorySuite) TestCountPreInputs() {
	field := cModel.APP_ID
	value := "1"
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	preCount, err := s.inputRepository.Count(s.ctx, filter)
	s.NoError(err)
	s.Equal(uint64(101), preCount)
}

func (s *InputRepositorySuite) TestCountInputs() {
	// Insert test data
	input1 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             model.InputCompletionStatus_Accepted,
	}
	input2 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             model.InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test counting all inputs
	count, err := s.inputRepository.Count(s.ctx, nil)
	s.NoError(err)
	s.Equal(uint64(103), count)

	field := "Status"
	value := model.InputCompletionStatus_Accepted.String()

	// Test counting inputs with specific status
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	count, err = s.inputRepository.Count(s.ctx, filter)
	s.NoError(err)
	s.Equal(uint64(100), count)
}

func (s *InputRepositorySuite) TestCountWrongStatusInputs() {
	field := cModel.STATUS_PROPERTY
	value := "CARTESI"

	// Test counting inputs with non-existent status
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	_, err := s.inputRepository.Count(s.ctx, filter)
	s.Error(err)
}

func (s *InputRepositorySuite) TestCountWrongAppIdInputs() {
	field := cModel.APP_ID
	value := "999" // non-existent application

	// Test counting inputs with non-existent status
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	count, err := s.inputRepository.Count(s.ctx, filter)
	s.NoError(err)
	s.Equal(uint64(0), count)
}

func (s *InputRepositorySuite) TestCountWrongFieldInputs() {
	field := cModel.APP_CONTRACT
	value := "0xdeadbeef" // non-existent application

	// Test counting inputs with non-existent status
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	_, err := s.inputRepository.Count(s.ctx, filter)
	s.Error(err)
}

func (s *InputRepositorySuite) TestFindAllInputsCount() {
	// Insert test data
	input1 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             model.InputCompletionStatus_Accepted,
	}
	input2 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             model.InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding all inputs
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, nil, nil, nil, nil)
	s.NoError(err)
	s.Len(inputs.Rows, 103)
}

func (s *InputRepositorySuite) TestFindAllInputsSpecificField() {
	// Insert test data
	input1 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             model.InputCompletionStatus_Accepted,
	}
	input2 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             model.InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding inputs with a specific status
	field := cModel.STATUS_PROPERTY
	value := model.InputCompletionStatus_Rejected.String()
	filter := []*cModel.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, nil, nil, nil, filter)
	s.NoError(err)
	s.Len(inputs.Rows, 1)
	s.Equal(input2.Index, inputs.Rows[0].Index)
}

func (s *InputRepositorySuite) TestFindAllInputsLimitOffset() {
	// Insert test data
	input1 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             model.InputCompletionStatus_Accepted,
	}
	input2 := model.Input{
		EpochApplicationID: 1,
		EpochIndex:         commons.OpenEpoch,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             model.InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding inputs with limit and offset
	last := 2
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, &last, nil, nil, nil)
	s.NoError(err)
	s.Len(inputs.Rows, 2)
	s.Equal(int(input1.Index), int(inputs.Rows[0].Index))
}

func TestInputRepositorySuite(t *testing.T) {
	suite.Run(t, new(InputRepositorySuite))
}
