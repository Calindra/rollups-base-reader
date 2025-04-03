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
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type AppRepositorySuite struct {
	suite.Suite
	appRepository *AppRepository
	ctx           context.Context
	ctxCancel     context.CancelFunc
	image         *postgres.PostgresContainer
	schemaDir     string
}

func TestAppRepository(t *testing.T) {
	suite.Run(t, new(AppRepositorySuite))
}

func (s *AppRepositorySuite) SetupSuite() {
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

func (s *AppRepositorySuite) SetupTest() {
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

	s.appRepository = NewAppRepository(db)
}

func (s *AppRepositorySuite) TearDownTest() {
	err := s.image.Stop(s.ctx, nil)
	s.NoError(err)
	s.appRepository.Db.Close()
	s.ctxCancel()
}

func (s *AppRepositorySuite) TestList() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	// Call FindAll
	apps, err := s.appRepository.List(ctx)
	s.NoError(err)
	s.Len(apps, 1)

	// Validate the first app
	firstApp := apps[0]
	s.Equal(int64(1), firstApp.ID)
	s.Equal("echo-dapp", firstApp.Name)
	s.Equal("0x8e3c7bF65833ccb1755dAB530Ef0405644FE6ae3", firstApp.IApplicationAddress.String())
}

func (s *AppRepositorySuite) TestFindOneByContract() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	contractAddress := common.HexToAddress("0x8e3c7bF65833ccb1755dAB530Ef0405644FE6ae3")

	// Call FindOne with application ID 1 (which is pre-populated in the test DB)
	app, err := s.appRepository.FindOneByContract(ctx, contractAddress)
	s.NoError(err)

	// Validate the application fields
	s.Equal(1, int(app.ID))
	s.Equal("echo-dapp", app.Name)
	s.Equal(contractAddress.Hex(), app.IApplicationAddress.String())
}
