// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputreader

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/contracts"
	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Model interface {
	AddAdvanceInput(
		sender common.Address,
		payload string,
		blockNumber uint64,
		timestamp time.Time,
		index int,
		prevRandao string,
		appContract common.Address,
		chainId string,
	) error
}

// This worker reads inputs from Ethereum and puts them in the model.
type InputReaderWorker struct {
	Model           Model
	Provider        string
	InputBoxAddress common.Address
	InputBoxBlock   uint64
	EthClient       *ethclient.Client
}

func (w InputReaderWorker) String() string {
	return "inputreader"
}

func (w *InputReaderWorker) GetEthClient() (*ethclient.Client, error) {
	if w.EthClient == nil {
		ctx := context.Background()
		client, err := ethclient.DialContext(ctx, w.Provider)
		if err != nil {
			return nil, fmt.Errorf("inputreader: dial: %w", err)
		}
		w.EthClient = client
	}
	return w.EthClient, nil
}

func (w *InputReaderWorker) ChainID() (*big.Int, error) {
	client, err := w.GetEthClient()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return client.ChainID(ctx)
}

func (w InputReaderWorker) FindAllInputsByBlockAndTimestampLT(
	ctx context.Context,
	client *ethclient.Client,
	inputBox *contracts.InputBox,
	startBlockNumber uint64,
	endTimestamp uint64,
	appAddresses []common.Address,
) ([]model.InputExtra, error) {
	slog.Debug("ReadInputsByBlockAndTimestamp",
		"startBlockNumber", startBlockNumber,
		"dappAddresses", appAddresses,
		"endTimestamp", endTimestamp,
	)

	opts := bind.FilterOpts{
		Context: ctx,
		Start:   startBlockNumber,
	}
	filter := appAddresses
	it, err := inputBox.FilterInputAdded(&opts, filter, nil)
	result := []model.InputExtra{}
	if err != nil {
		return result, fmt.Errorf("inputreader: filter input added: %v", err)
	}
	defer it.Close()

	for it.Next() {
		header, err := client.HeaderByHash(ctx, it.Event.Raw.BlockHash)
		if err != nil {
			return result, fmt.Errorf("inputreader: failed to get tx header: %w", err)
		}
		timestamp := uint64(header.Time)
		unixTimestamp := time.Unix(int64(header.Time), 0)
		if timestamp < endTimestamp {
			rawData := it.Event.Input
			eventInput := it.Event.Input[4:]
			abi, err := contracts.InputsMetaData.GetAbi()
			if err != nil {
				slog.Error("Error parsing abi", "err", err)
				return result, err
			}

			values, err := abi.Methods["EvmAdvance"].Inputs.UnpackValues(eventInput)
			if err != nil {
				slog.Error("Error parsing abi", "err", err)
				return result, err
			}

			chainId := values[0].(*big.Int).Uint64()
			appContract := values[1].(common.Address)
			msgSender := values[2].(common.Address)
			prevRandao := common.BytesToHash(values[5].(*big.Int).Bytes())
			payload := values[7].([]uint8)
			inputIndex := it.Event.Index.Uint64()

			input := model.InputExtra{
				Input: model.Input{
					Index:              inputIndex,
					BlockNumber:        header.Number.Uint64(),
					RawData:            rawData,
					Status:             model.InputCompletionStatus_None,
					EpochApplicationID: -1,
					EpochIndex:         0,
				},
				BlockTimestamp:  unixTimestamp,
				AppContract:     appContract,
				MsgSender:       msgSender,
				ChainId:         chainId,
				PrevRandao:      prevRandao,
				TransactionData: payload,
			}
			slog.Debug("append InputAdded", "timestamp", timestamp, "endTimestamp", endTimestamp)
			result = append(result, input)
		} else {
			slog.Debug("skip event InputAdded",
				"timestamp", timestamp,
				"endTimestamp", endTimestamp,
				"timeDiff", timestamp-endTimestamp,
			)
		}
	}

	return result, nil
}
