package repository

import (
	"context"
	"fmt"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/jmoiron/sqlx"
)

type AppRepository struct {
	Db *sqlx.DB
}

func NewAppRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{db}
}

func (a *AppRepository) FindAll(
	ctx context.Context,
	limit int,
	offset int,
	tx *sqlx.Tx,
) ([]model.Application, error) {
	query := `SELECT id, name, app_contract
	FROM application
	ORDER BY id
	LIMIT $1 OFFSET $2`
	args := []any{limit, offset}
	apps := []model.Application{}

	// Create a prepared statement
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing application query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.SelectContext(ctx, &apps, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying applications: %w", err)
	}

	return apps, nil
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