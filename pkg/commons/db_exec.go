package commons

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type staticstring string

const transactionKey staticstring = "safeWriteTransaction"

type DBExecutorInterface interface {
	sqlx.ExecerContext
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
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

	if isTxEnable {
		return tx.ExecContext(ctx, query, args...)
	} else {
		slog.Debug("Using ExecContext without transaction.")
		return d.db.ExecContext(ctx, query, args...)
	}
}

// PrepareNamedContext implements DBExecutorInterface.
func (d *DBExecutor) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	tx, isTxEnable := GetTransaction(ctx)

	if isTxEnable {
		return tx.PrepareNamedContext(ctx, query)
	} else {
		slog.Debug("Using PrepareNamedContext without transaction.")
		return d.db.PrepareNamedContext(ctx, query)
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
