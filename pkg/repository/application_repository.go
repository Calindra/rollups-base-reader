package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AppRepository struct {
	Db *sqlx.DB
}

func NewAppRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{db}
}

func (a *AppRepository) Exists(ctx context.Context, applicationId int64, tx *sqlx.Tx) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1
		FROM application
		WHERE id = $1
	)`
	args := []any{applicationId}
	var exists bool

	// Create a prepared statement
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("error preparing application exists query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.GetContext(ctx, &exists, args...)
	if err != nil {
		return false, fmt.Errorf("error checking if application exists: %w", err)
	}

	return exists, nil
}
