package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	TypeDuplicate        = "duplicate"
	TypeForeignKey       = "foreign_key"
	TypeNotNull          = "not_null"
	TypeCheckConstraint  = "check_constraint"
	TypeUndefinedObject  = "undefined_object"
	TypeAuthError        = "auth_error"
	TypeConnectionLimit  = "connection_limit"
	TypeDeadlock         = "deadlock"
	TypeQueryCanceled    = "query_canceled"
	TypePostgresError    = "postgres_error"
	TypeNoRows           = "no_rows"
	TypeTxClosed         = "tx_closed"
	TypeTxCommitRollback = "tx_commit_rollback"
	TypeTimeout          = "timeout"
	TypeConnectionError  = "connection_error"
	TypeUnknownError     = "unknown_error"
)

type PgxError struct {
	Type       string
	Original   error
	PgError    *pgconn.PgError
	Constraint string
	Table      string
	Column     string
	Detail     string
}

func ValidatePgxError(err error) *PgxError {
	if err == nil {
		return nil
	}

	result := &PgxError{Original: err}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		result.PgError = pgErr
		result.Detail = pgErr.Detail

		switch pgErr.Code {
		case "23505":
			result.Type = TypeDuplicate
			extractConstraintInfo(pgErr, result)
		case "23503":
			result.Type = TypeForeignKey
			extractConstraintInfo(pgErr, result)
		case "23502":
			result.Type = TypeNotNull
			extractConstraintInfo(pgErr, result)
		case "23514":
			result.Type = TypeCheckConstraint
			extractConstraintInfo(pgErr, result)
		case "42P01", "42704":
			result.Type = TypeUndefinedObject
		case "28000", "28P01":
			result.Type = TypeAuthError
		case "53300", "53400", "57P03":
			result.Type = TypeConnectionLimit
		case "40P01":
			result.Type = TypeDeadlock
		case "57014":
			result.Type = TypeQueryCanceled
		default:
			result.Type = TypePostgresError
		}
		return result
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		result.Type = TypeNoRows
	case errors.Is(err, pgx.ErrTxClosed):
		result.Type = TypeTxClosed
	case errors.Is(err, pgx.ErrTxCommitRollback):
		result.Type = TypeTxCommitRollback
	case isTimeoutError(err):
		result.Type = TypeTimeout
	case isConnectionError(err):
		result.Type = TypeConnectionError
	default:
		result.Type = TypeUnknownError
	}

	return result
}

func (p PgxError) String() string {
	return fmt.Sprintf("type: %s, original: %s, pgError: %s, constraint: %s, table: %s, column: %s, detail: %s",
		p.Type, p.Original, p.PgError, p.Constraint, p.Type, p.Column, p.Detail)
}

func extractConstraintInfo(pgErr *pgconn.PgError, result *PgxError) {
	parts := strings.Split(pgErr.Message, `"`)
	if len(parts) >= 2 {
		result.Constraint = parts[1]
	}

	if pgErr.Detail != "" {
		detailParts := strings.Split(pgErr.Detail, `(`)
		if len(detailParts) > 1 {
			tableAndColumn := strings.Trim(detailParts[1], `)`)
			tcParts := strings.Split(tableAndColumn, `, `)
			if len(tcParts) > 0 {
				result.Table = strings.Trim(tcParts[0], ` `)
			}
			if len(tcParts) > 1 {
				result.Column = strings.Trim(tcParts[1], ` `)
			}
		}
	}
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	return strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "timed out") ||
		strings.Contains(err.Error(), "context deadline exceeded")
}

func isConnectionError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && !netErr.Timeout() {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection") &&
		(strings.Contains(msg, "refused") ||
			strings.Contains(msg, "reset") ||
			strings.Contains(msg, "closed") ||
			strings.Contains(msg, "terminated"))
}
