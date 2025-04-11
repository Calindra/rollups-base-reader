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
	schemaPath    string
}

func TestAppRepository(t *testing.T) {
	suite.Run(t, new(AppRepositorySuite))
}

func (s *AppRepositorySuite) SetupSuite() {
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

func (s *AppRepositorySuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), commons.DefaultTimeout)

	// Database
	container, err := postgres.Run(s.ctx, commons.DbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaPath),
		postgres.WithDatabase(commons.DbName),
		postgres.WithUsername(commons.DbUser),
		postgres.WithPassword(commons.DbPassword),
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
	testcontainers.CleanupContainer(s.T(), s.image.Container)
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

func (s *AppRepositorySuite) TestUpdateDAByContract() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	secret := []byte("megasecret")[4:]
	da := model.DataAvailabilitySelector{secret[0], secret[1], secret[2], secret[3]}
	err := s.appRepository.UpdateDA(ctx, 1, da)
	s.NoError(err)
}

func (s *AppRepositorySuite) TestFindByDA() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	// Call FindByDA with InputBox DA selector
	apps, err := s.appRepository.FindAllByDA(ctx, model.DataAvailability_InputBox)
	s.NoError(err)

	// Verify we have application(s) using InputBox DA
	s.NotEmpty(apps, "Should find at least one application with InputBox data availability")

	// Check each app found has the correct DA selector
	for _, app := range apps {
		s.Equal(model.DataAvailability_InputBox, app.DataAvailability,
			"Application should have InputBox data availability")
	}

	// Create a custom DA selector for testing
	customDA := model.DataAvailabilitySelector{0xa1, 0xb2, 0xc3, 0xd4}

	// Query with the custom selector (which likely doesn't exist in the DB)
	customApps, err := s.appRepository.FindAllByDA(ctx, customDA)
	s.NoError(err)
	s.Empty(customApps, "Should not find any application with custom data availability")
}
