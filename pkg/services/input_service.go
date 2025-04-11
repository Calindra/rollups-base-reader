package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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
func (s *InputService) CreateInput(ctx context.Context, input model.Input) error {
	appID := input.EpochApplicationID

	// Get the latest open epoch for the app
	latestEpoch, err := s.EpochRepository.GetLatestOpenEpochByAppID(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to find latest epoch for appID %d: %w", appID, err)
	}

	if latestEpoch == nil {
		app, err := s.AppRepository.FindOneByID(ctx, appID)
		if err != nil {
			return fmt.Errorf("failed to find the app %d: %w", appID, err)
		}
		epoch := model.Epoch{
			ApplicationID: appID,
			Index:         0,
			FirstBlock:    input.BlockNumber,
			LastBlock:     input.BlockNumber + app.EpochLength,
			Status:        model.EpochStatus_Open,
			VirtualIndex:  0,
		}
		epochCreated, err := s.EpochRepository.Create(ctx, &epoch)
		if err != nil {
			return fmt.Errorf("failed to create an epoch for the app %d: %w", appID, err)
		}
		latestEpoch = epochCreated
		slog.Info("New epoch created")
	}

	// Set correct epoch index
	input.EpochIndex = latestEpoch.Index
	input.TransactionReference = &common.MaxHash

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

	return s.CreateInput(ctx, input)
}
