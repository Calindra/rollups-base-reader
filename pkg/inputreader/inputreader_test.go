// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputreader

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	util "github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/calindra/rollups-base-reader/pkg/contracts"
	"github.com/calindra/rollups-base-reader/pkg/devnet"
	"github.com/calindra/rollups-base-reader/pkg/repository"
	"github.com/cartesi/rollups-graphql/v2/pkg/commons"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type InputReaderTestSuite struct {
	suite.Suite
	appRepository repository.AppRepositoryInterface
	ctx           context.Context
	timeoutCancel context.CancelFunc
	schemaPath string
	postgresC  *postgres.PostgresContainer
	anvilC     *devnet.FoundryContainer
}

const timeout = 1 * time.Minute

func TestInputterTestSuite(t *testing.T) {
	suite.Run(t, new(InputReaderTestSuite))
}

func (s *InputReaderTestSuite) SetupSuite() {
	// Fetch schema
	tmpDir, err := os.MkdirTemp("", "schema")
	s.NoError(err)
	s.schemaPath = filepath.Join(tmpDir, "schema.sql")
	schemaFile, err := os.Create(s.schemaPath)
	s.NoError(err)
	defer schemaFile.Close()

	resp, err := http.Get(util.Schema)
	s.NoError(err)
	defer resp.Body.Close()

	_, err = io.Copy(schemaFile, resp.Body)
	s.NoError(err)
}

func (s *InputReaderTestSuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	slog.Debug("Setup!!!")

	s.ctx, s.timeoutCancel = context.WithTimeout(context.Background(), timeout)

	// Anvil
	anvilC, err := devnet.SetupFoundryV1(s.ctx)
	s.NoError(err)
	s.NotNil(anvilC)
	s.anvilC = anvilC

	// Database
	postgresC, err := postgres.Run(s.ctx, util.DbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaPath),
		postgres.WithDatabase(util.DbName),
		postgres.WithUsername(util.DbUser),
		postgres.WithPassword(util.DbPassword),
		testcontainers.WithLogConsumers(&util.StdoutLogConsumer{}),
	)
	s.NoError(err)
	extraArg := "sslmode=disable"
	connectionStr, err := postgresC.ConnectionString(s.ctx, extraArg)
	s.NoError(err)
	s.postgresC = postgresC
	err = postgresC.Start(s.ctx)
	s.NoError(err)

	db, err := sqlx.ConnectContext(s.ctx, "postgres", connectionStr)
	s.NoError(err)

	s.appRepository = repository.NewAppRepository(db)
}

func (s *InputReaderTestSuite) TearDownTest() {
	// Stop container
	testcontainers.CleanupContainer(s.T(), s.anvilC.Container)
	testcontainers.CleanupContainer(s.T(), s.postgresC.Container)

	s.appRepository.Close()
	s.timeoutCancel()
}

func (s *InputReaderTestSuite) TearDownSuite() {
	// Remove schema
	parentPath := filepath.Dir(s.schemaPath)
	err := os.RemoveAll(parentPath)
	s.NoError(err)
}

func (s *InputReaderTestSuite) TestFindAllInputsByBlockAndTimestampLT() {
	ctx, ctxCancel := context.WithCancel(s.ctx)
	defer ctxCancel()
	uri, err := s.anvilC.URI(ctx)
	s.NoError(err)
	client, err := ethclient.DialContext(ctx, uri)
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	err = devnet.AddInput(ctx, uri, common.Hex2Bytes("deadbeef"), devnet.ApplicationAddress)
	s.NoError(err)
	l1FinalizedPrevHeight := uint64(1)
	timestamp := uint64(time.Now().UnixMilli())
	w := InputReaderWorker{
		Model:           nil,
		Provider:        "",
		InputBoxAddress: inputBoxAddress,
		InputBoxBlock:   1,
	}

	inputs, err := w.FindAllInputsByBlockAndTimestampLT(ctx, client, inputBox, l1FinalizedPrevHeight, timestamp, []common.Address{appAddress})
	s.NoError(err)
	s.NotNil(inputs)
	s.Len(inputs, 1)
}

func (s *InputReaderTestSuite) TestZeroResultsFindAllInputsByBlockAndTimestampLT() {
	ctx, ctxCancel := context.WithCancel(s.ctx)
	defer ctxCancel()
	uri, err := s.anvilC.URI(ctx)
	s.NoError(err)
	client, err := ethclient.DialContext(ctx, uri)
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	err = devnet.AddInput(ctx, uri, common.Hex2Bytes("deadbeef"), devnet.ApplicationAddress)
	s.NoError(err)
	l1FinalizedPrevHeight := uint64(1)
	timestamp := uint64(time.Now().UnixMilli())
	w := InputReaderWorker{
		Model:           nil,
		Provider:        "",
		InputBoxAddress: inputBoxAddress,
		InputBoxBlock:   1,
	}
	inputs, err := w.FindAllInputsByBlockAndTimestampLT(ctx, client, inputBox, l1FinalizedPrevHeight, (timestamp/1000)-300, []common.Address{appAddress})
	s.NoError(err)
	s.NotNil(inputs)
	s.Len(inputs, 0)
}
