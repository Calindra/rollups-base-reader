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

type AppRepositoryInterface interface {
	FindOneByContract(ctx context.Context, address common.Address) (*model.Application, error)
	FindOneByID(ctx context.Context, id int64) (*model.Application, error)
	FindAllByDA(ctx context.Context, da model.DataAvailabilitySelector) ([]model.Application, error)
	UpdateDA(ctx context.Context, applicationId int64, da model.DataAvailabilitySelector) error
	List(ctx context.Context) ([]model.Application, error)
}

type AppRepository struct {
	Db *sqlx.DB
}

func NewAppRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{db}
}

// FindOneByContract returns a single application by contract address
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

// FindOneByID returns a single application by ID
func (a *AppRepository) FindOneByID(
	ctx context.Context,
	id int64,
) (*model.Application, error) {
	query := `SELECT *
	FROM application
	WHERE id = $1`
	args := []any{id}
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
		return nil, fmt.Errorf("error querying application with ID %d: %w", id, err)
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

	dbExec := commons.NewDBExecutor(a.Db)

	// Execute the query
	if _, err := dbExec.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("error updating application with ID %d: %w", applicationId, err)
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
