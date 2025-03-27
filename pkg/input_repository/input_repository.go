package inputrepository

import (
	"context"
	"fmt"

	nodev2 "github.com/cartesi/rollups-graphql/pkg/convenience/synchronizer_node"
	"github.com/jmoiron/sqlx"
)

type InputRepository struct {
	nodev2.RawRepository
}

func NewInputRepository(connectionURL string, db *sqlx.DB) *InputRepository {
	inside := nodev2.NewRawRepository(connectionURL, db)
	return &InputRepository{
		RawRepository: *inside,
	}
}

func (i *InputRepository) WriteInput(ctx context.Context, input Input) error {
	query := `INSERT INTO input (
		epoch_application_id,
		epoch_index,
		index,
		block_number,
		raw_data,
		status
	) VALUES ($1, $2, $3, $4, $5, $6)`
	args := []any{input.EpochApplicationID, input.EpochIndex, input.Index, input.BlockNumber, input.RawData, input.Status}

	// Create a prepared statement in Database
	stmt, err := i.Db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error preparing input query: %w", err)
	}
	defer stmt.Close()

	// Start a transaction
	tx, err := i.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	// Rollback the transaction if there is an error
	defer tx.Rollback()

	// Create a prepared statement in transaction
	stmtTx := tx.StmtContext(ctx, stmt)
	defer stmtTx.Close()

	// Execute the query
	_, err = stmtTx.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("error writing input: %w", err)
	}

	// Commit the transaction
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("error committing transaction: %w", commitErr)
	}

	return nil
}

func (i *InputRepository) QueryInput(ctx context.Context, applicationId int64, index uint64) (*Input, error) {
	query := `SELECT
		epoch_application_id,
		epoch_index,
		index,
		block_number,
		raw_data,
		status,
		created_at,
		updated_at
	FROM input
	WHERE epoch_application_id = $1 AND index = $2`
	args := []any{applicationId, index}

	// Create a prepared statement
	stmt, err := i.Db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	input := Input{}
	row := stmt.QueryRowContext(ctx, args...)
	if row == nil {
		return nil, fmt.Errorf("input not found")
	}
	err = row.Scan(&input.EpochApplicationID, &input.EpochIndex, &input.Index, &input.BlockNumber, &input.RawData, &input.Status, &input.CreatedAt, &input.UpdatedAt)

	return &input, err
}
