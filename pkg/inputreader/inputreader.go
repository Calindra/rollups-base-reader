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
	Model              Model
	Provider           string
	InputBoxAddress    common.Address
	InputBoxBlock      uint64
	ApplicationAddress common.Address
	EthClient          *ethclient.Client
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

// Read inputs starting from the input box deployment block until the latest block.
func (w *InputReaderWorker) ReadPastInputs(
	ctx context.Context,
	client *ethclient.Client,
	inputBox *contracts.InputBox,
	startBlockNumber uint64,
	endBlockNumber *uint64,
) error {
	if endBlockNumber != nil {
		slog.Debug("readPastInputs",
			"startBlockNumber", startBlockNumber,
			"endBlockNumber", *endBlockNumber,
			"dappAddress", w.ApplicationAddress,
		)
	} else {
		slog.Debug("readPastInputs",
			"startBlockNumber", startBlockNumber,
			"dappAddress", w.ApplicationAddress,
		)
	}
	opts := bind.FilterOpts{
		Context: ctx,
		Start:   startBlockNumber,
		End:     endBlockNumber,
	}
	filter := []common.Address{w.ApplicationAddress}
	it, err := inputBox.FilterInputAdded(&opts, filter, nil)
	if err != nil {
		return fmt.Errorf("inputreader: filter input added: %v", err)
	}
	defer it.Close()
	for it.Next() {
		w.InputBoxBlock = it.Event.Raw.BlockNumber - 1
		if err := w.addInput(ctx, client, it.Event); err != nil {
			return err
		}
	}
	return nil
}

// Add the input to the model.
func (w InputReaderWorker) addInput(
	ctx context.Context,
	client *ethclient.Client,
	event *contracts.InputBoxInputAdded,
) error {
	header, err := client.HeaderByHash(ctx, event.Raw.BlockHash)
	if err != nil {
		return fmt.Errorf("inputreader: failed to get tx header: %w", err)
	}
	timestamp := time.Unix(int64(header.Time), 0)

	// use abi to decode the input
	eventInput := event.Input[4:]
	abi, err := contracts.InputsMetaData.GetAbi()
	if err != nil {
		slog.Error("Error parsing abi", "err", err)
		return err
	}

	values, err := abi.Methods["EvmAdvance"].Inputs.UnpackValues(eventInput)
	if err != nil {
		slog.Error("Error parsing abi", "err", err)
		return err
	}

	chainId := values[0].(*big.Int).String()
	msgSender := values[2].(common.Address)
	prevRandao := fmt.Sprintf("0x%s", common.Bytes2Hex(values[5].(*big.Int).Bytes()))
	payload := common.Bytes2Hex(values[7].([]uint8))
	inputIndex := int(event.Index.Int64())

	slog.Debug("inputreader: read event",
		"dapp", event.AppContract,
		"input.index", event.Index,
		"sender", msgSender,
		"input", common.Bytes2Hex(event.Input),
		"payload", payload,
		slog.Group("block",
			"number", header.Number,
			"timestamp", timestamp,
			"prevRandao", prevRandao,
		),
	)

	if w.ApplicationAddress != event.AppContract {
		msg := fmt.Sprintf("The dapp address is wrong: %s. It should be %s",
			event.AppContract.Hex(),
			w.ApplicationAddress,
		)
		slog.Warn(msg)
		return nil
	}

	err = w.Model.AddAdvanceInput(
		msgSender,
		payload,
		event.Raw.BlockNumber,
		timestamp,
		inputIndex,
		prevRandao,
		event.AppContract,
		chainId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (w InputReaderWorker) FindAllInputsByBlockAndTimestampLT(
	ctx context.Context,
	client *ethclient.Client,
	inputBox *contracts.InputBox,
	startBlockNumber uint64,
	endTimestamp uint64,
) ([]model.Input, error) {
	slog.Debug("ReadInputsByBlockAndTimestamp",
		"startBlockNumber", startBlockNumber,
		"dappAddress", w.ApplicationAddress,
		"endTimestamp", endTimestamp,
	)

	opts := bind.FilterOpts{
		Context: ctx,
		Start:   startBlockNumber,
	}
	filter := []common.Address{w.ApplicationAddress}
	it, err := inputBox.FilterInputAdded(&opts, filter, nil)
	result := []model.Input{}
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
		// unixTimestamp := time.Unix(int64(header.Time), 0)
		if timestamp < endTimestamp {
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

			appContract := values[1].(common.Address)
			payload := values[7].([]uint8)
			inputIndex := it.Event.Index.Uint64()

			appStr := appContract.String()

			input := model.Input{
				Index:                inputIndex,
				BlockNumber:          header.Number.Uint64(),
				RawData:              payload,
				TransactionReference: it.Event.Raw.TxHash,
				Status:               model.InputCompletionStatus_None,
				EpochApplicationID:   -1,
				EpochIndex:           0,
				// a bit trick
				SnapshotURI: &appStr,
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
