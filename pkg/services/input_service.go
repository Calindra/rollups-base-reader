package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/calindra/rollups-base-reader/pkg/repository"
	"github.com/ethereum/go-ethereum/common"
)

type InputService struct {
	AppRepository   repository.AppRepositoryInterface
	InputRepository repository.InputRepositoryInterface
	EpochRepository repository.EpochRepositoryInterface
}

func NewInputService(inputRepository repository.InputRepositoryInterface, epochRepository repository.EpochRepositoryInterface, appRepository repository.AppRepositoryInterface) *InputService {
	return &InputService{
		InputRepository: inputRepository,
		EpochRepository: epochRepository,
		AppRepository:   appRepository,
	}
}
func (s *InputService) CreateInputID(ctx context.Context, appID int64, input model.Input) error {
	// Get the latest open epoch for the app
	latestEpoch, err := s.EpochRepository.GetLatestOpenEpochByAppID(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to find latest epoch for appID %d: %w", appID, err)
	}

	// Set correct epoch index
	input.EpochIndex = latestEpoch.Index

	err = s.InputRepository.Create(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create input: %w", err)
	}

	return nil
}

func (s *InputService) CreateInputWithAddress(ctx context.Context, appContract common.Address, input model.Input) error {
	// Check if the app exists
	app, err := s.AppRepository.FindOneByContract(ctx, appContract)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no app found: %w", err)
		}

		return fmt.Errorf("failed to find app: %w", err)
	}

	input.EpochApplicationID = app.ID

	return s.CreateInputID(ctx, app.ID, input)
}
