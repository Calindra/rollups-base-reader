// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputreader

import (
	"context"
	"fmt"
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
	"github.com/calindra/rollups-base-reader/pkg/supervisor"
	"github.com/cartesi/rollups-graphql/pkg/commons"
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
	appRepository *repository.AppRepository
	ctx           context.Context
	workerCtx     context.Context
	timeoutCancel context.CancelFunc
	workerCancel  context.CancelFunc
	workerResult  chan error
	rpcUrl        string
	schemaPath    string
	image         *postgres.PostgresContainer
}

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
	var w supervisor.SupervisorWorker
	w.Name = "TestInputter"
	s.ctx, s.timeoutCancel = context.WithTimeout(context.Background(), util.DefaultTimeout)
	s.workerResult = make(chan error)

	s.workerCtx, s.workerCancel = context.WithCancel(s.ctx)
	w.Workers = append(w.Workers, devnet.AnvilWorker{
		Address:  devnet.AnvilDefaultAddress,
		Port:     devnet.AnvilDefaultPort,
		Verbose:  true,
		AnvilCmd: "anvil",
	})

	// Database
	container, err := postgres.Run(s.ctx, util.DbImage,
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts(s.schemaPath),
		postgres.WithDatabase(util.DbName),
		postgres.WithUsername(util.DbUser),
		postgres.WithPassword(util.DbPassword),
		testcontainers.WithLogConsumers(&util.StdoutLogConsumer{}),
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

	s.appRepository = repository.NewAppRepository(db)

	s.rpcUrl = fmt.Sprintf("ws://%s:%v", devnet.AnvilDefaultAddress, devnet.AnvilDefaultPort)
	ready := make(chan struct{})
	go func() {
		s.workerResult <- w.Start(s.workerCtx, ready)
	}()
	select {
	case <-s.ctx.Done():
		s.Fail("context error", s.ctx.Err())
	case err := <-s.workerResult:
		s.Fail("worker exited before being ready", err)
	case <-ready:
		s.T().Log("nonodo ready")
	}
}

func (s *InputReaderTestSuite) TearDownTest() {
	s.workerCancel()
	select {
	case <-s.ctx.Done():
		s.Fail("context error", s.ctx.Err())
	case err := <-s.workerResult:
		s.NoError(err)
	}
	s.timeoutCancel()

	// Stop container
	testcontainers.CleanupContainer(s.T(), s.image.Container)

	s.appRepository.Db.Close()
}

func (s *InputReaderTestSuite) TearDownSuite() {
	// Remove schema
	parentPath := filepath.Dir(s.schemaPath)
	s.T().Logf("Removing schema path: %s", parentPath)
	err := os.RemoveAll(parentPath)
	s.NoError(err)
}

func (s *InputReaderTestSuite) TestFindAllInputsByBlockAndTimestampLT() {
	client, err := ethclient.DialContext(s.ctx, "http://127.0.0.1:8545")
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	ctx := context.Background()
	err = devnet.AddInput(ctx, s.rpcUrl, common.Hex2Bytes("deadbeef"), devnet.ApplicationAddress)
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
	client, err := ethclient.DialContext(s.ctx, "http://127.0.0.1:8545")
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	err = devnet.AddInput(ctx, s.rpcUrl, common.Hex2Bytes("deadbeef"), devnet.ApplicationAddress)
	s.NoError(err)
	l1FinalizedPrevHeight := uint64(1)
	timestamp := uint64(time.Now().UnixMilli())
	w := InputReaderWorker{
		Model:           nil,
		Provider:        "",
		InputBoxAddress: inputBoxAddress,
		InputBoxBlock:   1,
	}
	// block, err := client.BlockByNumber(ctx, nil)
	// s.NoError(err)
	// s.NotNil(block)
	// s.Equal(uint64(19), block.NumberU64())
	inputs, err := w.FindAllInputsByBlockAndTimestampLT(ctx, client, inputBox, l1FinalizedPrevHeight, (timestamp/1000)-300, []common.Address{appAddress})
	s.NoError(err)
	s.NotNil(inputs)
	s.Equal(0, len(inputs))
}
