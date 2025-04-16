package repository

import (
	"context"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OwnerRepositoryInterface interface {
	CreateOwner(ctx context.Context, tx pgx.Tx, data *dto.CreatePersonDTO, formID *uuid.UUID) error
}

type OwnerRepository struct {
}

func NewOwnerRepository() *OwnerRepository {
	return &OwnerRepository{}
}

func (o OwnerRepository) CreateOwner(
	ctx context.Context,
	tx pgx.Tx,
	data *dto.CreatePersonDTO,
	formID *uuid.UUID,
) error {
	query := `
		INSERT INTO 
		owner(email, telegram, name, surname, patronymic, relative, form_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.Exec(
		ctx,
		query,
		data.ContactEmail,
		data.ContactTelegram,
		data.ContactName,
		data.ContactSurname,
		data.ContactPatronymic,
		data.Relative,
		formID,
	); err != nil {
		return err
	}
	return nil
}
