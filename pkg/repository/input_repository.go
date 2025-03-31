package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cartesi/rollups-graphql/pkg/commons"
	"github.com/cartesi/rollups-graphql/pkg/convenience/model"
	cModel "github.com/cartesi/rollups-graphql/pkg/convenience/model"
	"github.com/cartesi/rollups-graphql/pkg/convenience/repository"
	"github.com/jmoiron/sqlx"
)

type staticstring string

const txKey staticstring = "safeTxWriteInput"

type InputRepository struct {
	Db              *sqlx.DB
	AppRepository   *AppRepository
	EpochRepository *EpochRepository
}

func NewInputRepository(db *sqlx.DB) *InputRepository {
	appRepo := NewAppRepository(db)
	epochRepo := NewEpochRepository(db)

	return &InputRepository{db, appRepo, epochRepo}
}

func (i *InputRepository) SafeCreate(stdCtx context.Context, input Input) error {
	tx, err := i.Db.BeginTxx(stdCtx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()
	ctx := context.WithValue(stdCtx, txKey, tx)

	exists, err := i.AppRepository.Exists(ctx, input.EpochApplicationID, tx)
	if err != nil {
		return fmt.Errorf("error checking if application exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("application with id %d does not exist", input.EpochApplicationID)
	}
	withinBounds, err := i.EpochRepository.isWithinBounds(ctx, input, tx)
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
			return i.Create(ctx, input)
		}
		return err
	}
	if inputDB != nil {
		return fmt.Errorf("input with index %d already exists", input.Index)
	}

	return err
}

func transformToInputQuery(
	filter []*model.ConvenienceFilter,
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
		case model.STATUS_PROPERTY:
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
		case model.APP_ID:
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
	filter []*model.ConvenienceFilter,
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

func (i *InputRepository) AdvanceInputToInput(advanceInput cModel.AdvanceInput) Input {
	panic("not implemented")
	// return Input{
	// 	EpochApplicationID: advanceInput.AppContract.Hex(),
	// 	EpochIndex:         advanceInput.Index,
	// 	Index:              advanceInput.Index,
	// 	BlockNumber:        advanceInput.BlockNumber,
	// 	RawData:            advanceInput.Payload,
	// 	Status:             advanceInput.Status,
	// }
}

func (i *InputRepository) Create(ctx context.Context, input Input) error {
	var (
		err error
		tx  *sqlx.Tx
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
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		// Start a transaction
		tx, err = i.Db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}
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
	stmt, err := i.Db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	input := &Input{}
	err = stmt.GetContext(ctx, input, args...)

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
	filter []*model.ConvenienceFilter,
) (*commons.PageResult[Input], error) {
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

	var inputs []Input
	err = stmt.SelectContext(ctx, &inputs, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing find all query: %w", err)
	}

	pageResult := &commons.PageResult[Input]{
		Rows:   inputs,
		Total:  total,
		Offset: uint64(offset),
	}

	return pageResult, nil
}
