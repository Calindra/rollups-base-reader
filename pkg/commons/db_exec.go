package commons

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type staticstring string

const txKey staticstring = "safeTxWriteInput"

func NewTx(ctx context.Context, db *sqlx.DB) (tx *sqlx.Tx, err error) {
	// Check if the transaction is already started
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		// Start a transaction
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("error starting transaction: %w", err)
		}
	}
	return tx, err
}
