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
	CreatePhoto(ctx context.Context, tx pgx.Tx, photo *dto.CreatePhotoDTO) error
	DeletePhoto(ctx context.Context, tx pgx.Tx, photoID int) error
}

type PhotoRepository struct {
}

func NewPhotoRepository() *PhotoRepository {
	return &PhotoRepository{}
}

func (PhotoRepository) CheckMainStatus(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (bool, error) {
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

func (PhotoRepository) CheckCount(ctx context.Context, tx pgx.Tx, countOK int, personID uuid.UUID) (bool, error) {
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

func (PhotoRepository) CreatePhoto(ctx context.Context, tx pgx.Tx, photo *dto.CreatePhotoDTO) error {
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

func (PhotoRepository) DeletePhoto(ctx context.Context, tx pgx.Tx, photoID int) error {
	query := `
	DELETE FROM person_photo
	WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, photoID)
	if err != nil {
		return err
	}

	return nil
}
