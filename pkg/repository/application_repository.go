package repository

import (
	"context"
	"fmt"

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
	query := `SELECT id, name, iapplication_address
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

func (a *AppRepository) List(
	ctx context.Context,
) ([]model.Application, error) {
	query := `SELECT id, name, iapplication_address
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
