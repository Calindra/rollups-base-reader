package commons

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type staticstring string

const transactionKey staticstring = "safeTxWriteInput"

type DBExecutorInterface interface {
	sqlx.ExecerContext
	sqlx.PreparerContext
}

type DBExecutor struct {
	db *sqlx.DB
}


func NewDBExecutor(db *sqlx.DB) DBExecutorInterface {
	return &DBExecutor{db: db}
}

// ExecContext implements sqlx.ExecerContext.
func (d *DBExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, isTxEnable := GetTransaction(ctx)

	if !isTxEnable {
		slog.Debug("Using ExecContext without transaction.")
		return d.db.ExecContext(ctx, query, args...)
	} else {
		return tx.ExecContext(ctx, query, args...)
	}
}

// PrepareContext implements DBExecutorInterface.
func (d *DBExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	tx, isTxEnable := GetTransaction(ctx)

	if !isTxEnable {
		slog.Debug("Using PrepareContext without transaction.")
		return d.db.PrepareContext(ctx, query)
	} else {
		return tx.PrepareContext(ctx, query)
	}
}

func StartTransaction(ctx context.Context, db *sqlx.DB) (context.Context, *sqlx.Tx, error) {
	tx, err := db.Beginx()
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	ctx = context.WithValue(ctx, transactionKey, tx)
	return ctx, tx, nil
}

func GetTransaction(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(transactionKey).(*sqlx.Tx)
	if !ok {
		slog.Debug("No transaction found in context")
		return nil, false
	}
	return tx, true
}
