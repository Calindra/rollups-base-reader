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

	_, err := i.Db.ExecContext(ctx, query, input.EpochApplicationID, input.EpochIndex, input.Index, input.BlockNumber, input.RawData, input.Status)
	if err != nil {
		return fmt.Errorf("error writing input: %w", err)
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

	input := Input{}
	row := i.Db.QueryRowContext(ctx, query, applicationId, index)
	if row == nil {
		return nil, fmt.Errorf("input not found")
	}
	err := row.Scan(&input.EpochApplicationID, &input.EpochIndex, &input.Index, &input.BlockNumber, &input.RawData, &input.Status, &input.CreatedAt, &input.UpdatedAt)
	return &input, err
}
