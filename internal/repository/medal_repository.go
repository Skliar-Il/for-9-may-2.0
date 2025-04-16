package repository

import (
	"context"
	"fmt"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MedalRepositoryInterface interface {
	CheckMedals(ctx context.Context, tx pgx.Tx, medals []int) (*bool, error)
	CreateMedalPerson(ctx context.Context, tx pgx.Tx, userID *uuid.UUID, medals []int) error
	GetMedals(ctx context.Context, tx pgx.Tx) ([]dto.MedalDTO, error)
	CreateMedal(ctx context.Context, tx pgx.Tx, medal *dto.CreateMedalDTO) (int, error)
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

func (m MedalRepository) GetMedals(ctx context.Context, tx pgx.Tx) ([]dto.MedalDTO, error) {
	query := `
		SELECT id, name, COALESCE(photo_link, '')
		FROM medal
		`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error(select medals): %w", err)
	}

	var medals []dto.MedalDTO
	for rows.Next() {
		var medal dto.MedalDTO

		if err := rows.Scan(
			&medal.ID,
			&medal.Name,
			&medal.ImageUrl,
		); err != nil {
			return nil, fmt.Errorf("scan medals rows error: %w", err)
		}
		medals = append(medals, medal)
	}
	return medals, nil
}

func (m MedalRepository) CreateMedal(ctx context.Context, tx pgx.Tx, medal *dto.CreateMedalDTO) (int, error) {
	query := `
		INSERT INTO medal(name, photo_link)
		VALUES ($1, $2)
		RETURNING id
		`

	var medalID int
	if err := tx.QueryRow(ctx, query, medal.Name, medal.ImageLink).Scan(&medalID); err != nil {
		return 0, err
	}
	return medalID, nil
}
