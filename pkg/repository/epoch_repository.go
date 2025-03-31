package repository

import (
	"context"
	"fmt"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/jmoiron/sqlx"
)

type EpochRepository struct {
	Db *sqlx.DB
}

func NewEpochRepository(db *sqlx.DB) *EpochRepository {
	return &EpochRepository{db}
}

func (e *EpochRepository) isWithinBounds(ctx context.Context, input model.Input, tx *sqlx.Tx) (bool, error) {
	epochIndex := input.EpochIndex
	query := `SELECT EXISTS (
		SELECT 1
		FROM epoch
		WHERE index >= $1 AND index <= (SELECT MAX(index) + 1 FROM epoch)
	)`
	args := []any{epochIndex}
	var withinBounds bool

	// Create a prepared statement
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("error preparing epoch bounds query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.GetContext(ctx, &withinBounds, args...)
	if err != nil {
		return false, fmt.Errorf("error checking if epoch is within bounds: %w", err)
	}

	return withinBounds, nil
}
