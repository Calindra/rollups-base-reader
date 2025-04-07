package services

import (
	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/calindra/rollups-base-reader/pkg/repository"
	cModel "github.com/cartesi/rollups-graphql/pkg/convenience/model"
)

type Convert struct {
	inputRepository *repository.InputRepository
}

func NewConvert(inputRepository *repository.InputRepository) *Convert {
	return &Convert{inputRepository}
}

func ConvertInputStatus(cInputStatus cModel.CompletionStatus) model.InputCompletionStatus {
	switch cInputStatus {
	case cModel.CompletionStatusAccepted:
		return model.InputCompletionStatus_Accepted
	case cModel.CompletionStatusRejected:
		return model.InputCompletionStatus_Rejected
	case cModel.CompletionStatusException:
		return model.InputCompletionStatus_Exception
	case cModel.CompletionStatusMachineHalted:
		return model.InputCompletionStatus_MachineHalted
	case cModel.CompletionStatusCycleLimitExceeded:
		return model.InputCompletionStatus_CycleLimitExceeded
	case cModel.CompletionStatusTimeLimitExceeded:
		return model.InputCompletionStatus_TimeLimitExceeded
	case cModel.CompletionStatusPayloadLengthLimitExceeded:
		return model.InputCompletionStatus_PayloadLengthLimitExceeded
	default:
		return model.InputCompletionStatus_None
	}
}
