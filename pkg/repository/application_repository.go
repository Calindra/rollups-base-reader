package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
)

type AppRepository struct {
	Db *sqlx.DB
}

func NewAppRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{db}
}

// FindOneByContract returns a single application by ID
func (a *AppRepository) FindOneByContract(
	ctx context.Context,
	address common.Address,
) (*model.Application, error) {
	query := `SELECT *
	FROM application
	WHERE iapplication_address = $1`
	args := []any{address}
	app := model.Application{}

	// Create a prepared statement
	stmt, err := a.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing application query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.GetContext(ctx, &app, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying application with address %s: %w", address.Hex(), err)
	}

	return &app, nil
}

func (a *AppRepository) FindAllByDA(ctx context.Context, da model.DataAvailabilitySelector) ([]model.Application, error) {
	query := `SELECT *
	FROM application
	WHERE data_availability = decode($1, 'hex')`
	daHex := common.Bytes2Hex(da[:])
	args := []any{daHex}
	apps := []model.Application{}

	slog.Debug("querying applications with data availability", "query", query, "args", args)
	// Create a prepared statement
	stmt, err := a.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing application query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.SelectContext(ctx, &apps, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying applications with data availability: %w", err)
	}

	return apps, nil
}

// Update field DA
func (a *AppRepository) UpdateDA(
	ctx context.Context,
	applicationId int64,
	da model.DataAvailabilitySelector,
) error {
	query := `UPDATE application
	SET data_availability = decode($1, 'hex')
	WHERE id = $2`
	daHex := common.Bytes2Hex(da[:])
	args := []any{daHex, applicationId}

	tx, err := commons.NewTx(ctx, a.Db)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Create a prepared statement
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error preparing application update query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("error updating application with ID %d: %w", applicationId, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	slog.Debug("updated application data availability", "applicationId", applicationId, "dataAvailability", da)
	return nil
}

func (a *AppRepository) List(
	ctx context.Context,
) ([]model.Application, error) {
	query := `SELECT *
	FROM application
	ORDER BY id`

	apps := []model.Application{}

	// Create a prepared statement
	stmt, err := a.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing application query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.SelectContext(ctx, &apps)
	if err != nil {
		return nil, fmt.Errorf("error querying applications: %w", err)
	}

	return apps, nil
}
