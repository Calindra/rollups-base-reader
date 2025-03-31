package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	util "github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/calindra/rollups-base-reader/pkg/model"
	"github.com/cartesi/rollups-graphql/pkg/commons"
	cModel "github.com/cartesi/rollups-graphql/pkg/convenience/model"
	"github.com/cartesi/rollups-graphql/pkg/convenience/repository"
	"github.com/jmoiron/sqlx"
)

type InputRepository struct {
	Db              *sqlx.DB
}

func NewInputRepository(db *sqlx.DB) *InputRepository {
	return &InputRepository{db}
}

func transformToInputQuery(
	filter []*cModel.ConvenienceFilter,
) (string, []any, int, error) {
	query := ""
	if len(filter) > 0 {
		query += repository.WHERE
	}
	args := []any{}
	where := []string{}
	count := 1
	for _, filter := range filter {
		switch *filter.Field {
		case repository.INDEX_FIELD:
			if filter.Eq != nil {
				where = append(where, fmt.Sprintf("index = $%d ", count))
				args = append(args, *filter.Eq)
				count += 1
			} else if filter.Gt != nil {
				where = append(where, fmt.Sprintf("index > $%d ", count))
				args = append(args, *filter.Gt)
				count += 1
			} else if filter.Lt != nil {
				where = append(where, fmt.Sprintf("index < $%d ", count))
				args = append(args, *filter.Lt)
				count += 1
			} else {
				return "", nil, 0, fmt.Errorf("operation not implemented")
			}
		case cModel.STATUS_PROPERTY:
			if filter.Ne != nil {
				where = append(where, fmt.Sprintf("status <> $%d ", count))
				args = append(args, *filter.Ne)
				count += 1
			} else if filter.Eq != nil {
				where = append(where, fmt.Sprintf("status = $%d ", count))
				args = append(args, *filter.Eq)
				count += 1
			} else {
				return "", nil, 0, fmt.Errorf("operation not implemented")
			}
		case cModel.APP_ID:
			if filter.Eq != nil {
				where = append(where, fmt.Sprintf("epoch_application_id = $%d ", count))
				args = append(args, *filter.Eq)
				count += 1
			} else {
				return "", nil, 0, fmt.Errorf("operation not implemented field epoch_application_id")
			}
		default:
			return "", nil, 0, fmt.Errorf("unexpected field %s", *filter.Field)
		}
	}
	query += strings.Join(where, " and ")
	return query, args, count, nil
}

func (i *InputRepository) Count(
	ctx context.Context,
	filter []*cModel.ConvenienceFilter,
) (uint64, error) {
	query := `SELECT count(*) FROM input `
	where, args, _, err := transformToInputQuery(filter)
	if err != nil {
		slog.Error("Count execution error", "err", err)
		return 0, err
	}
	query += where
	slog.Debug("Query", "query", query, "args", args)
	stmt, err := i.Db.PreparexContext(ctx, query)
	if err != nil {
		slog.Error("Count execution error")
		return 0, err
	}
	defer stmt.Close()
	var count uint64
	err = stmt.GetContext(ctx, &count, args...)
	if err != nil {
		slog.Error("Count execution error")
		return 0, err
	}
	return count, nil
}

func (i *InputRepository) Create(ctx context.Context, input model.Input) error {
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
	tx, err := util.NewTx(ctx, i.Db)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	// Rollback the transaction if there is an error
	defer tx.Rollback()

	// Create a prepared statement in transaction
	stmt, err := tx.PreparexContext(ctx, query)
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

func (i *InputRepository) FindOne(ctx context.Context, applicationId int64, index uint64) (*model.Input, error) {
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
	stmt, err := i.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	input := &model.Input{}
	row := stmt.QueryRowxContext(ctx, args...)
	err = row.Scan(&input.EpochApplicationID, &input.EpochIndex, &input.Index, &input.BlockNumber, &input.RawData, &input.Status, &input.CreatedAt, &input.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return input, nil
}

func (i *InputRepository) FindAll(
	ctx context.Context,
	first *int,
	last *int,
	after *string,
	before *string,
	filter []*cModel.ConvenienceFilter,
) (*commons.PageResult[model.Input], error) {
	total, err := i.Count(ctx, filter)
	if err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}
	query := `SELECT
		epoch_application_id,
		epoch_index,
		index,
		block_number,
		raw_data,
		status,
		created_at,
		updated_at
	FROM input `
	where, args, argsCount, err := transformToInputQuery(filter)
	if err != nil {
		return nil, fmt.Errorf("error transforming filter to query: %w", err)
	}
	query += where

	offset, limit, err := commons.ComputePage(first, last, after, before, int(total))
	if err != nil {
		return nil, err
	}
	query += fmt.Sprintf(`LIMIT $%d `, argsCount)
	args = append(args, limit)
	argsCount += 1
	query += fmt.Sprintf(`OFFSET $%d `, argsCount)
	args = append(args, offset)

	slog.Debug("Query", "query", query, "args", args, "total", total)
	stmt, err := i.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing find all query: %w", err)
	}
	defer stmt.Close()

	var inputs []model.Input
	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var input model.Input
		err = rows.Scan(&input.EpochApplicationID, &input.EpochIndex, &input.Index, &input.BlockNumber, &input.RawData, &input.Status, &input.CreatedAt, &input.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning input row: %w", err)
		}
		inputs = append(inputs, input)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	// err = stmt.SelectContext(ctx, &inputs, args...)
	// if err != nil {
	// 	return nil, fmt.Errorf("error executing find all query: %w", err)
	// }

	pageResult := &commons.PageResult[model.Input]{
		Rows:   inputs,
		Total:  total,
		Offset: uint64(offset),
	}

	return pageResult, nil
}
