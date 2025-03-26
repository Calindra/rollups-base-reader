package inputwritter

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type InputWritterSuite struct {
	suite.Suite
	inputWritter *InputRawWritter
	ctx          context.Context
	ctxCancel    context.CancelFunc
	image        *postgres.PostgresContainer
	schemaDir    string
}

const timeout time.Duration = 5 * time.Minute
const dbImage = "postgres:16-alpine"
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

func (s *InputWritterSuite) SetupSuite() {
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

func (s *InputWritterSuite) SetupTest() {
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

	s.inputWritter = NewInputRawWritter(connectionStr, db)
}

func (s *InputWritterSuite) TearDownTest() {
	err := s.image.Stop(s.ctx, nil)
	s.NoError(err)
	s.inputWritter.Db.Close()
	s.ctxCancel()
}

func (s *InputWritterSuite) TestWriteInput() {
	input := Input{
		EpochApplicationID: 1, // existing app
		EpochIndex:         23, // add to actual epoch
		Index:              171, // unique index
		BlockNumber:        0,
		RawData:            []byte("test data"),
		Status:             InputCompletionStatus_Accepted,
	}
	err := s.inputWritter.WriteInput(s.ctx, input)
	s.NoError(err)

	inputDb, err := s.inputWritter.QueryInput(s.ctx, input.EpochApplicationID, input.Index)
	s.NoError(err)
	s.Equal(input.EpochApplicationID, inputDb.EpochApplicationID)
}

func TestInputterWritterTestSuite(t *testing.T) {
	suite.Run(t, new(InputWritterSuite))
}
