package transaction

import (
	"context"

	"github.com/calindra/rollups-base-reader/pkg/eip712"
	"github.com/ethereum/go-ethereum/common"
)

type Sender interface {
	SubmitSigAndData(
		ctx context.Context,
		sigAndData eip712.SigAndData,
	) (string, error)
	GetNonce(
		ctx context.Context,
		appContract common.Address,
		msgSender common.Address,
	) (uint64, error)
}
