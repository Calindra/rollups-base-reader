package inputrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	nodev2 "github.com/cartesi/rollups-graphql/pkg/convenience/synchronizer_node"
	"github.com/jmoiron/sqlx"
)

type staticstring string

const txKey staticstring = "safeTxWriteInput"

type InputRepository struct {
	nodev2.RawRepository
}

func NewInputRepository(connectionURL string, db *sqlx.DB) *InputRepository {
	inside := nodev2.NewRawRepository(connectionURL, db)
	return &InputRepository{
		RawRepository: *inside,
	}
}

func (i *InputRepository) applicationExists(ctx context.Context, input Input, tx *sql.Tx) (bool, error) {
	applicationId := input.EpochApplicationID

	query := `SELECT EXISTS (
		SELECT 1
		FROM application
		WHERE id = $1
	)`
	args := []any{applicationId}
	var exists bool

	// Create a prepared statement
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("error preparing application exists query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRowContext(ctx, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if application exists: %w", err)
	}

	return exists, nil
}

func (i *InputRepository) isEpochWithinBounds(ctx context.Context, input Input, tx *sql.Tx) (bool, error) {
	epochIndex := input.EpochIndex
	query := `SELECT EXISTS (
		SELECT 1
		FROM epoch
		WHERE index >= $1 AND index <= (SELECT MAX(index) + 1 FROM epoch)
	)`
	var withinBounds bool

	// Create a prepared statement
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("error preparing epoch bounds query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRowContext(ctx, epochIndex).Scan(&withinBounds)
	if err != nil {
		return false, fmt.Errorf("error checking if epoch is within bounds: %w", err)
	}

	return withinBounds, nil
}

func (i *InputRepository) SafeWriteInput(stdCtx context.Context, input Input) error {
	tx, err := i.Db.BeginTx(stdCtx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()
	ctx := context.WithValue(stdCtx, txKey, tx)

	exists, err := i.applicationExists(ctx, input, tx)
	if err != nil {
		return fmt.Errorf("error checking if application exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("application with id %d does not exist", input.EpochApplicationID)
	}
	withinBounds, err := i.isEpochWithinBounds(ctx, input, tx)
	if err != nil {
		return fmt.Errorf("error checking if epoch is within bounds: %w", err)
	}
	if !withinBounds {
		return fmt.Errorf("epoch %d is not within bounds", input.EpochIndex)
	}

	// Check if the input already exists
	inputDB, err := i.QueryInput(ctx, input.EpochApplicationID, input.Index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return i.WriteInput(ctx, input)
		}
		return err
	}
	if inputDB != nil {
		return fmt.Errorf("input with index %d already exists", input.Index)
	}

	return err
}

func (i *InputRepository) WriteInput(ctx context.Context, input Input) error {
	var (
		err error
		tx  *sql.Tx
	)

	query := `INSERT INTO input (
		epoch_application_id,
		epoch_index,
		index,
		block_number,
		raw_data,
		status
	) VALUES ($1, $2, $3, $4, $5, $6)`
	args := []any{input.EpochApplicationID, input.EpochIndex, input.Index, input.BlockNumber, input.RawData, input.Status}

	// Check if the transaction is already started
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if !ok {
		// Start a transaction
		tx, err = i.Db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}
	}
	// Rollback the transaction if there is an error
	defer tx.Rollback()

	// Create a prepared statement in transaction
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error preparing input query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.ExecContext(ctx, args...)
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
	err = row.Scan(&input.EpochApplicationID, &input.EpochIndex, &input.Index, &input.BlockNumber, &input.RawData, &input.Status, &input.CreatedAt, &input.UpdatedAt)

	return &input, err
}
