// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputreader

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/contracts"
	"github.com/calindra/rollups-base-reader/pkg/devnet"
	"github.com/calindra/rollups-base-reader/pkg/supervisor"
	"github.com/cartesi/rollups-graphql/pkg/commons"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/suite"
)

type InputReaderTestSuite struct {
	suite.Suite
	ctx           context.Context
	workerCtx     context.Context
	timeoutCancel context.CancelFunc
	workerCancel  context.CancelFunc
	workerResult  chan error
	rpcUrl        string
}

func (s *InputReaderTestSuite) SetupTest() {
	commons.ConfigureLog(slog.LevelDebug)
	slog.Debug("Setup!!!")
	var w supervisor.SupervisorWorker
	w.Name = "TesteInputter"
	const testTimeout = 5 * time.Second
	s.ctx, s.timeoutCancel = context.WithTimeout(context.Background(), testTimeout)
	s.workerResult = make(chan error)

	s.workerCtx, s.workerCancel = context.WithCancel(s.ctx)
	// anvilLocation, err := devnet.CheckAnvilAndInstall(s.ctx)
	// s.NoError(err)
	w.Workers = append(w.Workers, devnet.AnvilWorker{
		Address:  devnet.AnvilDefaultAddress,
		Port:     devnet.AnvilDefaultPort,
		Verbose:  true,
		AnvilCmd: "anvil",
	})

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

func (s *InputReaderTestSuite) TestFindAllInputsByBlockAndTimestampLT() {
	client, err := ethclient.DialContext(s.ctx, "http://127.0.0.1:8545")
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	ctx := context.Background()
	err = devnet.AddInput(ctx, s.rpcUrl, common.Hex2Bytes("deadbeef"))
	s.NoError(err)
	l1FinalizedPrevHeight := uint64(1)
	timestamp := uint64(time.Now().UnixMilli())
	w := InputReaderWorker{
		Model:              nil,
		Provider:           "",
		InputBoxAddress:    inputBoxAddress,
		InputBoxBlock:      1,
		ApplicationAddress: appAddress,
	}

	inputs, err := w.FindAllInputsByBlockAndTimestampLT(ctx, client, inputBox, l1FinalizedPrevHeight, timestamp)
	s.NoError(err)
	s.NotNil(inputs)
	s.Equal(1, len(inputs))
}

func (s *InputReaderTestSuite) TestZeroResultsFindAllInputsByBlockAndTimestampLT() {
	client, err := ethclient.DialContext(s.ctx, "http://127.0.0.1:8545")
	s.NoError(err)
	appAddress := common.HexToAddress(devnet.ApplicationAddress)
	inputBoxAddress := common.HexToAddress(devnet.InputBoxAddress)
	inputBox, err := contracts.NewInputBox(inputBoxAddress, client)
	s.NoError(err)
	ctx := context.Background()
	err = devnet.AddInput(ctx, s.rpcUrl, common.Hex2Bytes("deadbeef"))
	s.NoError(err)
	l1FinalizedPrevHeight := uint64(1)
	timestamp := uint64(time.Now().UnixMilli())
	w := InputReaderWorker{
		Model:              nil,
		Provider:           "",
		InputBoxAddress:    inputBoxAddress,
		InputBoxBlock:      1,
		ApplicationAddress: appAddress,
	}
	// block, err := client.BlockByNumber(ctx, nil)
	// s.NoError(err)
	// s.NotNil(block)
	// s.Equal(uint64(19), block.NumberU64())
	inputs, err := w.FindAllInputsByBlockAndTimestampLT(ctx, client, inputBox, l1FinalizedPrevHeight, (timestamp/1000)-300)
	s.NoError(err)
	s.NotNil(inputs)
	s.Equal(0, len(inputs))
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
	s.T().Log("teardown ok.")
}

func TestInputterTestSuite(t *testing.T) {
	suite.Run(t, &InputReaderTestSuite{})
}
