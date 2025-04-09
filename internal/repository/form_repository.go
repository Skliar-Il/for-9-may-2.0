package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FormRepositoryInterface interface {
}

type FormRepository struct {
}

func (f FormRepository) CreateForm(ctx context.Context, tx pgx.Tx, person_id uuid.UUID) error {
	query := `
		INSERT INTO form(person_id)
		VALUES ($1)`
	if _, err := tx.Exec(ctx, query, person_id); err != nil {
		return err
	}
	return nil
}
