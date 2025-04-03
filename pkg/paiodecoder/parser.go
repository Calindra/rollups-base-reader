package paiodecoder

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os/exec"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/cartesi/rollups-graphql/pkg/commons"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const TimeoutExecutionPaioDecoder = 1 * time.Minute

type DecoderPaio interface {
	DecodePaioBatch(ctx context.Context, bytes []byte) (string, error)
}

type PaioDecoder struct {
	location string
}

func NewPaioDecoder(location string) *PaioDecoder {
	return &PaioDecoder{location}
}

// call the paio decoder binary
func (pd *PaioDecoder) DecodePaioBatch(stdCtx context.Context, rawBytes []byte) (string, error) {
	first, err := pd.DecodePaioBatchSkip(stdCtx, 0, rawBytes) // nolint
	if err == nil {
		return first, nil
	}
	slog.Warn("failed to decode, we will try again removing 2 bytes")
	second, err := pd.DecodePaioBatchSkip(stdCtx, 2, rawBytes) // nolint
	if err != nil {
		return "", err
	}
	return second, nil
}

func (pd *PaioDecoder) DecodePaioBatchSkip(stdCtx context.Context, skip int, rawBytes []byte) (string, error) {
	ctx, cancel := context.WithTimeout(stdCtx, TimeoutExecutionPaioDecoder)
	defer cancel()
	cmd := exec.CommandContext(ctx, pd.location)
	var stdinData bytes.Buffer
	bytesStr := common.Bytes2Hex(rawBytes[skip:])
	stdinData.WriteString(bytesStr)
	cmd.Stdin = &stdinData
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to decode", "bytes", bytesStr)
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	slog.Debug("Output decoded", "output", string(output))
	return string(output), nil
}

func CreateTypedData(
	app common.Address,
	nonce uint64,
	maxGasPrice *big.Int,
	dataBytes []byte,
	chainId *big.Int,
) apitypes.TypedData {
	var typedData apitypes.TypedData
	cid := math.NewHexOrDecimal256(chainId.Int64())
	typedData.Domain = commons.NewCartesiDomain(cid)
	typedData.Types = apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"CartesiMessage": {
			{Name: "app", Type: "address"},
			{Name: "nonce", Type: "uint64"},
			{Name: "max_gas_price", Type: "uint128"},
			{Name: "data", Type: "bytes"},
		},
	}
	typedData.PrimaryType = "CartesiMessage"
	typedData.Message = apitypes.TypedDataMessage{
		"app":           app.String(),
		"nonce":         nonce,
		"max_gas_price": maxGasPrice.String(),
		"data":          fmt.Sprintf("0x%s", common.Bytes2Hex(dataBytes)),
	}
	return typedData
}

func ParsePaioBatchToInputs(jsonStr string, chainId *big.Int) ([]model.InputExtra, error) {
	inputs := []model.InputExtra{}
	var paioBatch PaioBatch
	if err := json.Unmarshal([]byte(jsonStr), &paioBatch); err != nil {
		return inputs, fmt.Errorf("unmarshal paio batch: %w", err)
	}
	slog.Debug("PaioBatch", "tx len", len(paioBatch.Txs), "json", jsonStr)
	for _, tx := range paioBatch.Txs {
		slog.Debug("Tx",
			"app", tx.App,
			"signature", tx.Signature.Hex(),
		)
		typedData := CreateTypedData(
			common.HexToAddress(tx.App),
			tx.Nonce,
			big.NewInt(int64(tx.MaxGasPrice)),
			tx.Data,
			chainId,
		)
		typeJSON, err := json.Marshal(typedData)
		if err != nil {
			return inputs, fmt.Errorf("error marshalling typed data: %w", err)
		}
		// set the typedData as string json below
		sigAndData := commons.SigAndData{
			Signature: tx.Signature.Hex(),
			TypedData: base64.StdEncoding.EncodeToString(typeJSON),
		}
		jsonPayload, err := json.Marshal(sigAndData)
		if err != nil {
			slog.Error("Error json.Marshal message:", "err", err)
			return inputs, err
		}
		slog.Debug("SaveTransaction", "jsonPayload", string(jsonPayload))
		msgSender, _, signature, err := commons.ExtractSigAndData(string(jsonPayload))
		if err != nil {
			slog.Error("Error ExtractSigAndData message:", "err", err)
			return inputs, err
		}

		txHex := fmt.Sprintf("0x%s", common.Bytes2Hex(crypto.Keccak256(signature)))
		data := AvailData{
			Transaction: txHex,
			Data:        tx.Data,
			AppContract: common.HexToAddress(tx.App),
			MsgSender:   msgSender,
		}

		encoded, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling data: %w", err)
		}

		slog.Debug("Tx", "hash", txHex)

		input := model.InputExtra{
			MsgSender:   msgSender,
			AppContract: common.HexToAddress(tx.App),
			ChainId:     chainId.Uint64(),
			Input: model.Input{
				Index:              0,
				BlockNumber:        0,
				Status:             model.InputCompletionStatus_None,
				EpochIndex:         0,
				EpochApplicationID: 0,
				RawData:            encoded,
			},
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func ParsePaioFrom712Message(typedData apitypes.TypedData) (PaioMessage, error) {
	message := PaioMessage{
		App:         typedData.Message["app"].(string),
		Nonce:       typedData.Message["nonce"].(string),
		MaxGasPrice: typedData.Message["max_gas_price"].(string),
		Payload:     []byte(typedData.Message["data"].(string)),
	}
	return message, nil
}
