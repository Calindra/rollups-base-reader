package repository

// Inspired from https://github.com/cartesi/rollups-node/blob/c748b205fe35f217868d98e8e64909afdd11d9af/internal/model/models.go

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type InputCompletionStatus string

const (
	InputCompletionStatus_None                       InputCompletionStatus = "NONE"
	InputCompletionStatus_Accepted                   InputCompletionStatus = "ACCEPTED"
	InputCompletionStatus_Rejected                   InputCompletionStatus = "REJECTED"
	InputCompletionStatus_Exception                  InputCompletionStatus = "EXCEPTION"
	InputCompletionStatus_MachineHalted              InputCompletionStatus = "MACHINE_HALTED"
	InputCompletionStatus_OutputsLimitExceeded       InputCompletionStatus = "OUTPUTS_LIMIT_EXCEEDED"
	InputCompletionStatus_CycleLimitExceeded         InputCompletionStatus = "CYCLE_LIMIT_EXCEEDED"
	InputCompletionStatus_TimeLimitExceeded          InputCompletionStatus = "TIME_LIMIT_EXCEEDED"
	InputCompletionStatus_PayloadLengthLimitExceeded InputCompletionStatus = "PAYLOAD_LENGTH_LIMIT_EXCEEDED"
)

var InputCompletionStatusAllValues = []InputCompletionStatus{
	InputCompletionStatus_None,
	InputCompletionStatus_Accepted,
	InputCompletionStatus_Rejected,
	InputCompletionStatus_Exception,
	InputCompletionStatus_MachineHalted,
	InputCompletionStatus_CycleLimitExceeded,
	InputCompletionStatus_TimeLimitExceeded,
	InputCompletionStatus_PayloadLengthLimitExceeded,
}

func (e *InputCompletionStatus) Scan(value any) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("invalid value for InputCompletionStatus enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "NONE":
		*e = InputCompletionStatus_None
	case "ACCEPTED":
		*e = InputCompletionStatus_Accepted
	case "REJECTED":
		*e = InputCompletionStatus_Rejected
	case "EXCEPTION":
		*e = InputCompletionStatus_Exception
	case "MACHINE_HALTED":
		*e = InputCompletionStatus_MachineHalted
	case "CYCLE_LIMIT_EXCEEDED":
		*e = InputCompletionStatus_CycleLimitExceeded
	case "TIME_LIMIT_EXCEEDED":
		*e = InputCompletionStatus_TimeLimitExceeded
	case "PAYLOAD_LENGTH_LIMIT_EXCEEDED":
		*e = InputCompletionStatus_PayloadLengthLimitExceeded
	default:
		return errors.New("Invalid value '" + enumValue + "' for InputCompletionStatus enum")
	}

	return nil
}

func (e InputCompletionStatus) String() string {
	return string(e)
}

type Input struct {
	EpochApplicationID   int64                 `db:"epoch_application_id" sql:"primary_key" json:"-"`
	EpochIndex           uint64                `db:"epoch_index" json:"epoch_index"`
	Index                uint64                `db:"index" sql:"primary_key" json:"index"`
	BlockNumber          uint64                `db:"block_number" json:"block_number"`
	RawData              []byte                `db:"raw_data" json:"raw_data"`
	Status               InputCompletionStatus `db:"status" json:"status"`
	MachineHash          *common.Hash          `db:"machine_hash" json:"machine_hash"`
	OutputsHash          *common.Hash          `db:"outputs_hash" json:"outputs_hash"`
	TransactionReference common.Hash           `db:"transaction_reference" json:"transaction_reference"`
	SnapshotURI          *string               `db:"snapshot_uri" json:"-"`
	CreatedAt            time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time             `db:"updated_at" json:"updated_at"`
}
