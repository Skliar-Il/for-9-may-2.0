package database

import (
	"context"
	"errors"
	"fmt"
	"for9may/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func RollbackTx(ctx context.Context, tx pgx.Tx, logger *logger.Logger) {
	if tx != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			if !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				logger.Error(ctx, fmt.Sprintf("failed complite transaction: %v", rollbackErr))
			}
		}
	}
}
