// Source: https://raw.githubusercontent.com/cartesi/rollups-node/d8063929ef309df6925ff45883cce482ec3d9e18/internal/model/models.go

// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package model

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Application struct {
	ID                   int64                    `db:"id" json:"-"`
	Name                 string                   `db:"name" json:"name"`
	IApplicationAddress  common.Address           `db:"iapplication_address" json:"iapplication_address"`
	IConsensusAddress    common.Address           `db:"iconsensus_address" json:"iconsensus_address"`
	IInputBoxAddress     common.Address           `db:"iinputbox_address" json:"iinputbox_address"`
	TemplateHash         common.Hash              `db:"template_hash" json:"template_hash"`
	TemplateURI          string                   `db:"template_uri" json:"-"`
	EpochLength          uint64                   `db:"epoch_length" json:"epoch_length"`
	DataAvailability     DataAvailabilitySelector `db:"data_availability" json:"data_availability"`
	State                ApplicationState         `db:"state" json:"state"`
	Reason               *string                  `db:"reason" json:"reason"`
	IInputBoxBlock       uint64                   `db:"iinputbox_block" json:"iinputbox_block"`
	LastInputCheckBlock  uint64                   `db:"last_input_check_block" json:"last_input_check_block"`
	LastOutputCheckBlock uint64                   `db:"last_output_check_block" json:"last_output_check_block"`
	ProcessedInputs      uint64                   `db:"processed_inputs" json:"processed_inputs"`
	CreatedAt            time.Time                `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time                `db:"updated_at" json:"updated_at"`
	ExecutionParameters  ExecutionParameters      `db:"execution_parameters" json:"execution_parameters"`
}

func (a *Application) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion in MarshalJSON.
	type Alias Application
	// Define a new structure that embeds the alias but overrides the hex fields.
	aux := &struct {
		*Alias
		IInputBoxBlock       string `json:"iinputbox_block"`
		LastInputCheckBlock  string `json:"last_input_check_block"`
		LastOutputCheckBlock string `json:"last_output_check_block"`
	}{
		Alias:                (*Alias)(a),
		IInputBoxBlock:       fmt.Sprintf("0x%x", a.IInputBoxBlock),
		LastInputCheckBlock:  fmt.Sprintf("0x%x", a.LastInputCheckBlock),
		LastOutputCheckBlock: fmt.Sprintf("0x%x", a.LastOutputCheckBlock),
	}
	return json.Marshal(aux)
}

type ApplicationState string

const (
	ApplicationState_Enabled    ApplicationState = "ENABLED"
	ApplicationState_Disabled   ApplicationState = "DISABLED"
	ApplicationState_Inoperable ApplicationState = "INOPERABLE"
)

var ApplicationStateAllValues = []ApplicationState{
	ApplicationState_Enabled,
	ApplicationState_Disabled,
	ApplicationState_Inoperable,
}

func (e *ApplicationState) Scan(value any) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("Invalid value for ApplicationState enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "ENABLED":
		*e = ApplicationState_Enabled
	case "DISABLED":
		*e = ApplicationState_Disabled
	case "INOPERABLE":
		*e = ApplicationState_Inoperable
	default:
		return errors.New("Invalid value '" + enumValue + "' for ApplicationState enum")
	}

	return nil
}

func (e ApplicationState) String() string {
	return string(e)
}

const DATA_AVAILABILITY_SELECTOR_SIZE = 4

type DataAvailabilitySelector [DATA_AVAILABILITY_SELECTOR_SIZE]byte

// Known data availability selectors
var (
	// ABI encoded "InputBox(address)"
	DataAvailability_InputBox = DataAvailabilitySelector{0xb1, 0x2c, 0x9e, 0xde}
)

func (d *DataAvailabilitySelector) MarshalJSON() ([]byte, error) {
	return json.Marshal("0x" + hex.EncodeToString(d[:]))
}

func (d *DataAvailabilitySelector) Scan(value any) error {
	var selector []byte
	switch v := value.(type) {
	case []byte:
		selector = v
	default:
		return errors.New("Invalid scan value for DataAvailabilitySelector. Value has to be of type []byte")
	}

	if len(selector) != DATA_AVAILABILITY_SELECTOR_SIZE {
		return errors.New("Invalid value for DataAvailabilitySelector")
	}
	copy(d[:], selector[:DATA_AVAILABILITY_SELECTOR_SIZE])

	return nil
}

type SnapshotPolicy string

const (
	SnapshotPolicy_None      SnapshotPolicy = "NONE"
	SnapshotPolicy_EachInput SnapshotPolicy = "EACH_INPUT"
	SnapshotPolicy_EachEpoch SnapshotPolicy = "EACH_EPOCH"
)

var SnapshotPolicyAllValues = []SnapshotPolicy{
	SnapshotPolicy_None,
	SnapshotPolicy_EachInput,
	SnapshotPolicy_EachEpoch,
}

func (e *SnapshotPolicy) Scan(value any) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("Invalid scan value for SnapshotPolicy enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "NONE":
		*e = SnapshotPolicy_None
	case "EACH_INPUT":
		*e = SnapshotPolicy_EachInput
	case "EACH_EPOCH":
		*e = SnapshotPolicy_EachEpoch
	default:
		return errors.New("Invalid scan value '" + enumValue + "' for SnapshotPolicy enum")
	}

	return nil
}

func (e SnapshotPolicy) String() string {
	return string(e)
}

type ExecutionParameters struct {
	ApplicationID         int64          `sql:"primary_key" json:"-"`
	SnapshotPolicy        SnapshotPolicy `json:"snapshot_policy"`
	SnapshotRetention     uint64         `json:"snapshot_retention"`
	AdvanceIncCycles      uint64         `json:"advance_inc_cycles"`
	AdvanceMaxCycles      uint64         `json:"advance_max_cycles"`
	InspectIncCycles      uint64         `json:"inspect_inc_cycles"`
	InspectMaxCycles      uint64         `json:"inspect_max_cycles"`
	AdvanceIncDeadline    time.Duration  `json:"advance_inc_deadline"`
	AdvanceMaxDeadline    time.Duration  `json:"advance_max_deadline"`
	InspectIncDeadline    time.Duration  `json:"inspect_inc_deadline"`
	InspectMaxDeadline    time.Duration  `json:"inspect_max_deadline"`
	LoadDeadline          time.Duration  `json:"load_deadline"`
	StoreDeadline         time.Duration  `json:"store_deadline"`
	FastDeadline          time.Duration  `json:"fast_deadline"`
	MaxConcurrentInspects uint32         `json:"max_concurrent_inspects"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

type Epoch struct {
	ApplicationID        int64        `sql:"primary_key" db:"application_id" json:"-"`
	Index                uint64       `sql:"primary_key" db:"index" json:"index"`
	FirstBlock           uint64       `db:"first_block" json:"first_block"`
	LastBlock            uint64       `db:"last_block" json:"last_block"`
	ClaimHash            *common.Hash `db:"claim_hash" json:"claim_hash"`
	ClaimTransactionHash *common.Hash `db:"claim_transaction_hash" json:"claim_transaction_hash"`
	Status               EpochStatus  `db:"status" json:"status"`
	VirtualIndex         uint64       `db:"virtual_index" json:"virtual_index"`
	CreatedAt            time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time    `db:"updated_at" json:"updated_at"`
}

func (e *Epoch) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion in MarshalJSON.
	type Alias Epoch
	// Define a new structure that embeds the alias but overrides the hex fields.
	aux := &struct {
		Index        string `json:"index"`
		FirstBlock   string `json:"first_block"`
		LastBlock    string `json:"last_block"`
		VirtualIndex string `json:"virtual_index"`
		*Alias
	}{
		Index:        fmt.Sprintf("0x%x", e.Index),
		FirstBlock:   fmt.Sprintf("0x%x", e.FirstBlock),
		LastBlock:    fmt.Sprintf("0x%x", e.LastBlock),
		VirtualIndex: fmt.Sprintf("0x%x", e.VirtualIndex),
		Alias:        (*Alias)(e),
	}
	return json.Marshal(aux)
}

type EpochStatus string

const (
	EpochStatus_Open            EpochStatus = "OPEN"
	EpochStatus_Closed          EpochStatus = "CLOSED"
	EpochStatus_InputsProcessed EpochStatus = "INPUTS_PROCESSED"
	EpochStatus_ClaimComputed   EpochStatus = "CLAIM_COMPUTED"
	EpochStatus_ClaimSubmitted  EpochStatus = "CLAIM_SUBMITTED"
	EpochStatus_ClaimAccepted   EpochStatus = "CLAIM_ACCEPTED"
	EpochStatus_ClaimRejected   EpochStatus = "CLAIM_REJECTED"
)

var EpochStatusAllValues = []EpochStatus{
	EpochStatus_Open,
	EpochStatus_Closed,
	EpochStatus_InputsProcessed,
	EpochStatus_ClaimComputed,
	EpochStatus_ClaimSubmitted,
	EpochStatus_ClaimAccepted,
	EpochStatus_ClaimRejected,
}

func (e *EpochStatus) Scan(value any) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("Invalid value for EpochStatus enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "OPEN":
		*e = EpochStatus_Open
	case "CLOSED":
		*e = EpochStatus_Closed
	case "INPUTS_PROCESSED":
		*e = EpochStatus_InputsProcessed
	case "CLAIM_COMPUTED":
		*e = EpochStatus_ClaimComputed
	case "CLAIM_SUBMITTED":
		*e = EpochStatus_ClaimSubmitted
	case "CLAIM_ACCEPTED":
		*e = EpochStatus_ClaimAccepted
	case "CLAIM_REJECTED":
		*e = EpochStatus_ClaimRejected
	default:
		return errors.New("Invalid value '" + enumValue + "' for EpochStatus enum")
	}

	return nil
}

func (e EpochStatus) String() string {
	return string(e)
}

type Input struct {
	EpochApplicationID   int64                 `sql:"primary_key" db:"epoch_application_id" json:"-"`
	EpochIndex           uint64                `db:"epoch_index" json:"epoch_index"`
	Index                uint64                `sql:"primary_key" db:"index" json:"index"`
	BlockNumber          uint64                `db:"block_number" json:"block_number"`
	RawData              []byte                `db:"raw_data" json:"raw_data"`
	Status               InputCompletionStatus `db:"status" json:"status"`
	MachineHash          *common.Hash          `db:"machine_hash" json:"machine_hash"`
	OutputsHash          *common.Hash          `db:"outputs_hash" json:"outputs_hash"`
	TransactionReference *common.Hash          `db:"transaction_reference" json:"transaction_reference"`
	SnapshotURI          *string               `db:"snapshot_uri" json:"-"`
	CreatedAt            time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time             `db:"updated_at" json:"updated_at"`
}

func (i *Input) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion in MarshalJSON.
	type Alias Input
	// Define a new structure that embeds the alias but overrides the hex fields.
	aux := &struct {
		EpochIndex  string `json:"epoch_index"`
		Index       string `json:"index"`
		BlockNumber string `json:"block_number"`
		RawData     string `json:"raw_data"`
		*Alias
	}{
		EpochIndex:  fmt.Sprintf("0x%x", i.EpochIndex),
		Index:       fmt.Sprintf("0x%x", i.Index),
		BlockNumber: fmt.Sprintf("0x%x", i.BlockNumber),
		RawData:     "0x" + hex.EncodeToString(i.RawData),
		Alias:       (*Alias)(i),
	}
	return json.Marshal(aux)
}

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
		return errors.New("Invalid value for InputCompletionStatus enum. Enum value has to be of type string or []byte")
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

type Output struct {
	InputEpochApplicationID  int64         `sql:"primary_key" json:"-"`
	InputIndex               uint64        `json:"input_index"`
	Index                    uint64        `sql:"primary_key" json:"index"`
	RawData                  []byte        `json:"raw_data"`
	Hash                     *common.Hash  `json:"hash"`
	OutputHashesSiblings     []common.Hash `json:"output_hashes_siblings"`
	ExecutionTransactionHash *common.Hash  `json:"execution_transaction_hash"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

func (i *Output) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion in MarshalJSON.
	type Alias Output
	// Define a new structure that embeds the alias but overrides the hex fields.
	aux := &struct {
		InputIndex string `json:"input_index"`
		Index      string `json:"index"`
		RawData    string `json:"raw_data"`
		*Alias
	}{
		InputIndex: fmt.Sprintf("0x%x", i.InputIndex),
		Index:      fmt.Sprintf("0x%x", i.Index),
		RawData:    "0x" + hex.EncodeToString(i.RawData),
		Alias:      (*Alias)(i),
	}
	return json.Marshal(aux)
}

type Report struct {
	InputEpochApplicationID int64     `sql:"primary_key" json:"-"`
	InputIndex              uint64    `json:"input_index"`
	Index                   uint64    `sql:"primary_key" json:"index"`
	RawData                 []byte    `json:"raw_data"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (r *Report) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion in MarshalJSON.
	type Alias Report
	// Define a new structure that embeds the alias but overrides the hex fields.
	aux := &struct {
		InputIndex string `json:"input_index"`
		Index      string `json:"index"`
		RawData    string `json:"raw_data"`
		*Alias
	}{
		InputIndex: fmt.Sprintf("0x%x", r.InputIndex),
		Index:      fmt.Sprintf("0x%x", r.Index),
		RawData:    "0x" + hex.EncodeToString(r.RawData),
		Alias:      (*Alias)(r),
	}
	return json.Marshal(aux)
}

type NodeConfig[T any] struct {
	Key       string
	Value     T
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AdvanceResult struct {
	InputIndex  uint64
	Status      InputCompletionStatus
	Outputs     [][]byte
	Reports     [][]byte
	OutputsHash common.Hash
	MachineHash *common.Hash
}

type InspectResult struct {
	ProcessedInputs uint64
	Accepted        bool
	Reports         [][]byte
	Error           error
}

// FIXME: remove this type. Migrate claim to use Application + Epoch
type ClaimRow struct {
	Epoch
	IApplicationAddress common.Address
	IConsensusAddress   common.Address
}

type DefaultBlock string

const (
	DefaultBlock_Finalized DefaultBlock = "FINALIZED"
	DefaultBlock_Latest    DefaultBlock = "LATEST"
	DefaultBlock_Pending   DefaultBlock = "PENDING"
	DefaultBlock_Safe      DefaultBlock = "SAFE"
)

var DefaultBlockAllValues = []DefaultBlock{
	DefaultBlock_Finalized,
	DefaultBlock_Latest,
	DefaultBlock_Pending,
	DefaultBlock_Safe,
}

func (e *DefaultBlock) Scan(value any) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("Invalid value for DefaultBlock enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "FINALIZED":
		*e = DefaultBlock_Finalized
	case "LATEST":
		*e = DefaultBlock_Latest
	case "PENDING":
		*e = DefaultBlock_Pending
	case "SAFE":
		*e = DefaultBlock_Safe
	default:
		return errors.New("Invalid value '" + enumValue + "' for DefaultBlock enum")
	}

	return nil
}

func (e DefaultBlock) String() string {
	return string(e)
}

type MonitoredEvent string

const (
	MonitoredEvent_InputAdded     MonitoredEvent = "InputAdded"
	MonitoredEvent_OutputExecuted MonitoredEvent = "OutputExecuted"
	MonitoredEvent_ClaimSubmitted MonitoredEvent = "ClaimSubmitted"
	MonitoredEvent_ClaimAccepted  MonitoredEvent = "ClaimAccepted"
)

func (e MonitoredEvent) String() string {
	return string(e)
}

func Pointer[T any](v T) *T {
	return &v
}

// Extra
type InputExtra struct {
	Input
	BlockTimestamp time.Time
	AppContract    common.Address
	MsgSender      common.Address
	ChainId        uint64
	PrevRandao     common.Hash
}
