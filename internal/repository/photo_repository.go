package repository

import (
	"context"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PhotoRepositoryInterface interface {
	CheckMainStatus(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (bool, error)
	CheckCount(ctx context.Context, tx pgx.Tx, countOK int, personID uuid.UUID) (bool, error)
	CreatePhoto(ctx context.Context, tx pgx.Tx, photo *dto.CreateNewPhotoDTO) error
}

type PhotoRepository struct {
}

func NewPhotoRepository() *PersonRepository {
	return &PersonRepository{}
}

func (PersonRepository) CheckMainStatus(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM person_photo
		WHERE person_id = $1
	`
	var count int
	if err := tx.QueryRow(ctx, query, personID).Scan(&count); err != nil {
		return false, err
	}

	if count >= 1 {
		return true, nil
	}

	return false, nil
}

func (PersonRepository) CheckCount(ctx context.Context, tx pgx.Tx, countOK int, personID uuid.UUID) (bool, error) {
	query := `
	SELECT COUNT(*)
	FROM person_photo
	WHERE person_id = $1
	`
	var count int
	if err := tx.QueryRow(ctx, query, personID).Scan(&count); err != nil {
		return false, err
	}

	if countOK <= count {
		return false, nil
	}

	return true, nil
}

func (PersonRepository) CreatePhoto(ctx context.Context, tx pgx.Tx, photo *dto.CreateNewPhotoDTO) error {
	query := `
		INSERT INTO person_photo(person_id, link, main_status)
		VALUES($1, $2, $3)
	`
	_, err := tx.Exec(ctx, query, photo.PersonID, photo.Link, photo.MainStatus)
	if err != nil {
		return err
	}

	return nil
}
