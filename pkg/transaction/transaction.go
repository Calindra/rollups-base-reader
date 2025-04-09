package transaction

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi.yaml oapi-transaction.yaml

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/calindra/rollups-base-reader/pkg/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
)

type TransactionAPI struct {
	ClientSender Sender
}

var _ ServerInterface = (*TransactionAPI)(nil)

// SaveTransaction implements ServerInterface.
func (p *TransactionAPI) SaveTransaction(ctx echo.Context) error {
	panic("unimplemented")
}

// SendTransaction implements ServerInterface.
func (p *TransactionAPI) SendTransaction(ctx echo.Context) error {
	panic("unimplemented")
}

// SendCartesiTransaction implements ServerInterface.
func (p *TransactionAPI) SendCartesiTransaction(ctx echo.Context) error {
	stdCtx, cancel := context.WithCancel(ctx.Request().Context())
	defer cancel()
	var request SendCartesiTransactionJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	typeJSON, err := json.Marshal(request.TypedData)
	if err != nil {
		return fmt.Errorf("error marshalling typed data: %w", err)
	}
	sigAndData := eip712.SigAndData{
		Signature: *request.Signature,
		TypedData: base64.StdEncoding.EncodeToString(typeJSON),
	}
	jsonPayload, err := json.Marshal(sigAndData)
	if err != nil {
		slog.Error("Error json.Marshal message:", "err", err)
		return err
	}
	slog.Debug("/submit", "jsonPayload", string(jsonPayload))
	msgSender, _, signature, err := eip712.ExtractSigAndData(string(jsonPayload))
	if err != nil {
		slog.Error("Error ExtractSigAndData message:", "err", err)
		return err
	}
	if request.Address != nil && common.HexToAddress(*request.Address) != msgSender {
		errorMessage := "wrong signature"
		return ctx.JSON(http.StatusBadRequest, TransactionError{Message: &errorMessage})
	}
	appContract := common.HexToAddress(request.TypedData.Message.App[2:])
	slog.Debug("SaveTransaction",
		"msgSender", msgSender,
		"appContract", appContract.Hex(),
		"message", request.TypedData.Message,
	)
	txId := fmt.Sprintf("0x%s", common.Bytes2Hex(crypto.Keccak256(signature)))

	seqTxId, err := p.ClientSender.SubmitSigAndData(stdCtx, sigAndData)
	if err != nil {
		return err
	}
	slog.Info("Transaction sent to the sequencer", "txId", seqTxId)
	response := TransactionResponse{
		Id: &txId,
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (p *TransactionAPI) GetNonce(ctx echo.Context) error {
	var request GetNonceJSONRequestBody
	stdCtx, cancel := context.WithCancel(ctx.Request().Context())
	defer cancel()
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if request.MsgSender == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "msg_sender is required"})
	}
	msgSender := common.HexToAddress(request.MsgSender)
	appContract := common.HexToAddress(request.AppContract)

	slog.Debug("GetNonce", "AppContract", request.AppContract, "MsgSender", request.MsgSender)

	total, err := p.ClientSender.GetNonce(stdCtx, appContract, msgSender)
	if err != nil {
		slog.Error("Error querying for inputs:", "err", err)
		return err
	}
	nonce := int(total)
	response := NonceResponse{
		Nonce: &nonce,
	}

	return ctx.JSON(http.StatusOK, response)
}
