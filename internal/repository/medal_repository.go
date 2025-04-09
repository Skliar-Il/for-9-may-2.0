package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MedalRepositoryInterface interface {
	CheckMedals(ctx context.Context, tx pgx.Tx, medals []int) (*bool, error)
	CreateMedalPerson(ctx context.Context, tx pgx.Tx, userID *uuid.UUID, medals []int) error
}
type MedalRepository struct{}

func NewMedalRepository() *MedalRepository {
	return &MedalRepository{}
}

func (m MedalRepository) CheckMedals(ctx context.Context, tx pgx.Tx, medals []int) (*bool, error) {
	query := `
		SELECT COUNT(DISTINCT id) = $1 AS all_exist
		FROM medal
		WHERE id = ANY($2)`

	var status bool
	err := tx.QueryRow(ctx, query, len(medals), medals).Scan(&status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (m MedalRepository) CreateMedalPerson(ctx context.Context, tx pgx.Tx, userID *uuid.UUID, medals []int) error {
	query := `
		INSERT INTO medal_person(person_id, medal_id)
		VALUES ($1, $2)`

	for _, medalID := range medals {
		if _, err := tx.Exec(ctx, query, userID, medalID); err != nil {
			return err
		}
	}
	return nil
}
