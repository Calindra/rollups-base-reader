package repository

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type EpochRepositorySuite struct {
	suite.Suite
	epochRepository *EpochRepository
	ctx             context.Context
	ctxCancel       context.CancelFunc
	image           *postgres.PostgresContainer
	schemaPath       string
}

func TestEpochRepository(t *testing.T) {
	suite.Run(t, new(EpochRepositorySuite))
}

func (s *EpochRepositorySuite) SetupSuite() {
	// Fetch schema
	tmpDir, err := os.MkdirTemp("", "schema")
	s.NoError(err)
	s.schemaPath = filepath.Join(tmpDir, "schema.sql")
	schemaFile, err := os.Create(s.schemaPath)
	s.NoError(err)
	defer schemaFile.Close()

	resp, err := http.Get(commons.Schema)
	s.NoError(err)
	defer resp.Body.Close()

	_, err = io.Copy(schemaFile, resp.Body)
	s.NoError(err)
}

func (s *EpochRepositorySuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), commons.DefaultTimeout)

	// Database
	container, err := postgres.Run(s.ctx, commons.DbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaPath),
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

	s.epochRepository = NewEpochRepository(db)
}

func (s *EpochRepositorySuite) TearDownTest() {
	testcontainers.CleanupContainer(s.T(), s.image.Container)
	s.epochRepository.Db.Close()
	s.ctxCancel()
}

func (s *EpochRepositorySuite) TestGetLatestOpenEpoch() {
	ctx, ctxCancel := context.WithCancel(s.ctx)
	defer ctxCancel()

	epoch, err := s.epochRepository.GetLatestOpenEpoch(ctx)
	s.NoError(err)
	s.NotNil(epoch)
	s.Equal(19, int(epoch.Index))
	s.Equal(model.EpochStatus_Open, epoch.Status)
}

func (s *EpochRepositorySuite) TestFindOne() {
	ctx, ctxCancel := context.WithCancel(s.ctx)
	defer ctxCancel()

	epoch, err := s.epochRepository.FindOne(ctx, 18)
	s.NoError(err)
	s.NotNil(epoch)
	s.Equal(18, int(epoch.Index))
	s.Equal(model.EpochStatus_ClaimComputed, epoch.Status)
}
