package inputrepository

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/cartesi/rollups-graphql/pkg/convenience/model"
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

const timeout time.Duration = 5 * time.Minute
const dbImage = "postgres:17.4-alpine"
const dbName = "rollups"
const dbUser = "postgres"
const dbPassword = "password"

const schema = "https://raw.githubusercontent.com/cartesi/rollups-graphql/8e63682b2e99282cdfee3f13b608d0316c22a484/postgres/raw/rollupsdb-dump-202503191059.sql"

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	if l.LogType == testcontainers.StderrLog {
		slog.Error(string(l.Content))
		return
	}

	slog.Debug(string(l.Content))
}

func (s *InputRepositorySuite) SetupSuite() {
	// Fetch schema
	tmpDir, err := os.MkdirTemp("", "schema")
	s.NoError(err)
	s.schemaDir = filepath.Join(tmpDir, "schema.sql")
	schemaFile, err := os.Create(s.schemaDir)
	s.NoError(err)
	defer schemaFile.Close()

	resp, err := http.Get(schema)
	s.NoError(err)
	defer resp.Body.Close()

	_, err = io.Copy(schemaFile, resp.Body)
	s.NoError(err)
}

func (s *InputRepositorySuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), timeout)

	// Database
	container, err := postgres.Run(s.ctx, dbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaDir),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithLogConsumers(&StdoutLogConsumer{}),
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
	input := Input{
		EpochApplicationID: 1,   // existing app
		EpochIndex:         23,  // add to actual epoch
		Index:              171, // unique index
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}
	err := s.inputRepository.Create(s.ctx, input)
	s.NoError(err)

	inputDb, err := s.inputRepository.QueryInput(s.ctx, input.EpochApplicationID, input.Index)
	s.NoError(err)
	s.Equal(input.EpochApplicationID, inputDb.EpochApplicationID)
}

func (s *InputRepositorySuite) TestInputWrongIndex() {
	input := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              1,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestSafeInputWrongIndex() {
	input := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              1,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.SafeCreate(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestInputWrongEpoch() {
	input := Input{
		EpochApplicationID: 1,
		EpochIndex:         999, // non-existent epoch
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestSafeInputWrongEpoch() {
	input := Input{
		EpochApplicationID: 1,
		EpochIndex:         999, // non-existent epoch
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.SafeCreate(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestInputWrongApplication() {
	input := Input{
		EpochApplicationID: 999, // non-existent application
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.Create(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestSafeInputWrongApplication() {
	input := Input{
		EpochApplicationID: 999, // non-existent application
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}

	err := s.inputRepository.SafeCreate(s.ctx, input)
	s.Error(err)
}

func (s *InputRepositorySuite) TestQueryInputWrongApplicationIndex() {
	var (
		applicationId int64  = 999 // non-existent application
		index         uint64 = 171 // non-existent index
	)
	_, err := s.inputRepository.QueryInput(s.ctx, applicationId, index)
	s.Error(err)
	s.ErrorIs(err, sql.ErrNoRows)
}

func (s *InputRepositorySuite) TestCountPreInputs() {
	field := model.APP_ID
	value := "1"
	filter := []*model.ConvenienceFilter{
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
	input1 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             InputCompletionStatus_Accepted,
	}
	input2 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             InputCompletionStatus_Rejected,
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
	value := InputCompletionStatus_Accepted.String()

	// Test counting inputs with specific status
	filter := []*model.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	count, err = s.inputRepository.Count(s.ctx, filter)
	s.NoError(err)
	s.Equal(uint64(102), count)
}

func (s *InputRepositorySuite) TestCountWrongStatusInputs() {
	field := model.STATUS_PROPERTY
	value := "CARTESI"

	// Test counting inputs with non-existent status
	filter := []*model.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	_, err := s.inputRepository.Count(s.ctx, filter)
	s.Error(err)
}

func (s *InputRepositorySuite) TestCountWrongAppIdInputs() {
	field := model.APP_ID
	value := "999" // non-existent application

	// Test counting inputs with non-existent status
	filter := []*model.ConvenienceFilter{
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
	field := model.APP_CONTRACT
	value := "0xdeadbeef" // non-existent application

	// Test counting inputs with non-existent status
	filter := []*model.ConvenienceFilter{
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
	input1 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             InputCompletionStatus_Accepted,
	}
	input2 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding all inputs
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, nil, nil, nil, nil)
	s.NoError(err)
	s.Len(inputs, 103)
}

func (s *InputRepositorySuite) TestFindAllInputsSpecificField() {
	// Insert test data
	input1 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             InputCompletionStatus_Accepted,
	}
	input2 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding inputs with a specific status
	field := model.STATUS_PROPERTY
	value := InputCompletionStatus_Rejected.String()
	filter := []*model.ConvenienceFilter{
		{
			Field: &field,
			Eq:    &value,
		},
	}
	inputs, err := s.inputRepository.FindAll(s.ctx, filter, nil, nil, nil, nil)
	s.NoError(err)
	s.Len(inputs, 1)
	s.Equal(input2.Index, inputs[0].Index)
}

func (s *InputRepositorySuite) TestFindAllInputsLimitOffset() {
	// Insert test data
	input1 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             InputCompletionStatus_Accepted,
	}
	input2 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding inputs with limit and offset
	limit := uint64(1)
	offset := uint64(101)
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, &limit, &offset, nil, nil)
	s.NoError(err)
	s.Len(inputs, 1)
	s.Equal(input1.Index, inputs[0].Index)
}

func (s *InputRepositorySuite) TestFindAllInputsOrderDirection() {
	// Insert test data
	input1 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              171,
		BlockNumber:        0,
		RawData:            []byte("test data 1"),
		Status:             InputCompletionStatus_Accepted,
	}
	input2 := Input{
		EpochApplicationID: 1,
		EpochIndex:         23,
		Index:              172,
		BlockNumber:        0,
		RawData:            []byte("test data 2"),
		Status:             InputCompletionStatus_Rejected,
	}
	err := s.inputRepository.Create(s.ctx, input1)
	s.NoError(err)
	err = s.inputRepository.Create(s.ctx, input2)
	s.NoError(err)

	// Test finding inputs with order by and direction
	orderBy := "index"
	orderDirection := "DESC"
	inputs, err := s.inputRepository.FindAll(s.ctx, nil, nil, nil, &orderBy, &orderDirection)
	s.NoError(err)
	s.Len(inputs, 103)
	s.Equal(input2.Index, inputs[0].Index)
	s.Equal(input1.Index, inputs[1].Index)
}

func TestInputRepositorySuite(t *testing.T) {
	suite.Run(t, new(InputRepositorySuite))
}
