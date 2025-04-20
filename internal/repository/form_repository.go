package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FormRepositoryInterface interface {
	CreateForm(ctx context.Context, tx pgx.Tx, personID *uuid.UUID) (*uuid.UUID, error)
	StatusForm(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (bool, error)
	UpdateForm(ctx context.Context, tx pgx.Tx, personID uuid.UUID, status bool) error
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

func (f FormRepository) StatusForm(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (bool, error) {
	query := `
		SELECT status_check
		FROM form
		WHERE person_id = $1
	`
	var formStatus bool
	if err := tx.QueryRow(ctx, query, personID).Scan(&formStatus); err != nil {
		return false, err
	}

	return formStatus, nil
}

func (FormRepository) UpdateForm(ctx context.Context, tx pgx.Tx, personID uuid.UUID, status bool) error {
	query := `
		UPDATE form
		SET
			main = $1
		WHERE person_id = $2
	`
	_, err := tx.Exec(ctx, query, status, personID)
	return err
}
