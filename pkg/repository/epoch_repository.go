package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/jmoiron/sqlx"
)

type EpochRepositoryInterface interface {
	GetLatestOpenEpochByAppID(ctx context.Context, appID int64) (*model.Epoch, error)
	FindOne(ctx context.Context, index uint64) (*model.Epoch, error)
	Create(ctx context.Context, epoch *model.Epoch) (*model.Epoch, error)
}

type EpochRepository struct {
	Db *sqlx.DB
}

func NewEpochRepository(db *sqlx.DB) *EpochRepository {
	return &EpochRepository{db}
}

func (e *EpochRepository) GetLatestOpenEpochByAppID(ctx context.Context, appID int64) (*model.Epoch, error) {
	query := `SELECT
		index,
		first_block,
		last_block,
		claim_hash,
		claim_transaction_hash,
		status,
		virtual_index,
		created_at,
		updated_at
	FROM epoch
	WHERE status = $1 AND application_id = $2
	ORDER BY index DESC
	LIMIT 1`
	args := []any{model.EpochStatus_Open, appID}

	epoch := model.Epoch{}

	// Create a prepared statement
	stmt, err := e.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.GetContext(ctx, &epoch, args...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &epoch, nil
}

// FindOne retrieves a specific epoch by its index
func (e *EpochRepository) FindOne(ctx context.Context, index uint64) (*model.Epoch, error) {
	query := `SELECT
		index,
		first_block,
		last_block,
		claim_hash,
		claim_transaction_hash,
		status,
		virtual_index,
		created_at,
		updated_at
	FROM epoch
	WHERE index = $1`
	args := []any{index}

	epoch := model.Epoch{}

	// Create a prepared statement
	stmt, err := e.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.GetContext(ctx, &epoch, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil to indicate no record found
		}
		return nil, err
	}

	return &epoch, nil
}

func (e *EpochRepository) Create(ctx context.Context, epoch *model.Epoch) (*model.Epoch, error) {
	query := `
		INSERT INTO epoch (
			application_id,
			index,
			first_block,
			last_block,
			claim_hash,
			claim_transaction_hash,
			status,
			virtual_index
		) VALUES (
			:application_id,
			:index,
			:first_block,
			:last_block,
			:claim_hash,
			:claim_transaction_hash,
			:status,
			:virtual_index
		)
		RETURNING 
			application_id,
			index,
			first_block,
			last_block,
			claim_hash,
			claim_transaction_hash,
			status,
			virtual_index,
			created_at,
			updated_at
	`

	stmt, err := e.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var inserted model.Epoch
	err = stmt.GetContext(ctx, &inserted, epoch)
	if err != nil {
		return nil, err
	}

	return &inserted, nil
}
