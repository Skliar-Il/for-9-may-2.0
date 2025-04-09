package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FormRepositoryInterface interface {
	CreateForm(ctx context.Context, tx pgx.Tx, personID *uuid.UUID) (*uuid.UUID, error)
}

type FormRepository struct {
}

func NewFormRepository() *FormRepository {
	return &FormRepository{}
}

func (f FormRepository) CreateForm(ctx context.Context, tx pgx.Tx, personID *uuid.UUID) (*uuid.UUID, error) {
	query := `
		INSERT INTO form(person_id)
		VALUES ($1)
		RETURNING id`
	var formID uuid.UUID
	if err := tx.QueryRow(ctx, query, personID).Scan(&formID); err != nil {
		return nil, err
	}
	return &formID, nil
}
