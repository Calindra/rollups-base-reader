// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package devnet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type AnvilSuite struct {
	suite.Suite
	ctx       context.Context
	ctxCancel context.CancelFunc
}

//
// Suite entry point
//

func TestAnvilSuite(t *testing.T) {
	suite.Run(t, &AnvilSuite{})
}

func (s *AnvilSuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), testTimeout)
}

func (s *AnvilSuite) TearDownTest() {
	s.ctxCancel()
}

const testTimeout = 30 * time.Second

func (s *AnvilSuite) TestAnvilWorker() {
	ctx, timeoutCancel := context.WithCancel(s.ctx)
	defer timeoutCancel()

	anvilCmd, err := CheckAnvilAndInstall(ctx)
	s.Require().NoError(err)

	anvilPort := AnvilDefaultPort + 100
	w := AnvilWorker{
		Address:  AnvilDefaultAddress,
		Port:     anvilPort,
		Verbose:  true,
		AnvilCmd: anvilCmd,
	}

	// start worker in goroutine
	workerCtx, workerCancel := context.WithCancel(ctx)
	defer workerCancel()
	ready := make(chan struct{})
	result := make(chan error)
	go func() {
		result <- w.Start(workerCtx, ready)
	}()

	// wait until worker is ready
	select {
	case <-ready:
	case <-ctx.Done():
		s.NoError(ctx.Err())
	}

	// send input
	rpcUrl := fmt.Sprintf("http://127.0.0.1:%v", anvilPort)
	payload := common.Hex2Bytes("deadbeef")
	err = AddInput(ctx, rpcUrl, payload, ApplicationAddress)
	s.NoError(err)

	// read input
	events, err := GetInputAdded(ctx, rpcUrl)
	s.NoError(err)
	s.Equal(1, len(events))

	// check input
	abi, err := contracts.InputsMetaData.GetAbi()
	s.NoError(err)

	values, err := abi.Methods["EvmAdvance"].Inputs.UnpackValues(events[0].Input[4:])
	s.NoError(err)

	receivedPayload := values[7].([]byte)

	s.Equal(payload, receivedPayload)

	// stop worker
	workerCancel()
	canceled := false
	select {
	case err := <-result:
		s.Equal(context.Canceled, err)
		canceled = true
	case <-ctx.Done():
		s.NoError(ctx.Err())
	}
	s.True(canceled)
}

func (s *AnvilSuite) TestGetContract() {
	contracts, err := GetContractInfo()
	s.NoError(err)
	s.NotEmpty(contracts)
}

type addressBook struct {
	*bytes.Buffer
}

// Close implements io.WriteCloser.
func (a *addressBook) Close() error {
	return nil
}

var _ io.WriteCloser = (*addressBook)(nil)

func (s *AnvilSuite) TestAddressbookContract() {
	output := &addressBook{bytes.NewBufferString("")}
	err := ShowAddresses(output)
	s.NoError(err)
	s.NotEmpty(output.String())
	slog.Info("output", "output", output.String())
}
