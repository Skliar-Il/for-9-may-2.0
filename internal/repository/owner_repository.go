package repository

import (
	"context"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OwnerRepositoryInterface interface {
	Create(ctx context.Context, tx pgx.Tx, data *dto.CreatePersonDTO, formID *uuid.UUID) error
	Update(ctx context.Context, tx pgx.Tx, person *dto.UpdatePersonDTO) error
}

type OwnerRepository struct {
}

func NewOwnerRepository() *OwnerRepository {
	return &OwnerRepository{}
}

func (OwnerRepository) Create(
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

func (OwnerRepository) Update(
	ctx context.Context,
	tx pgx.Tx,
	person *dto.UpdatePersonDTO,
) error {
	query := `
		UPDATE owner
		SET
			email = $1,
			telegram = $2,
			name = $3,
			surname = $4,
			patronymic = $5,
			relative = $6
	`

	_, err := tx.Exec(ctx, query, person.ContactEmail, person.ContactTelegram, person.ContactName,
		person.ContactSurname, person.ContactPatronymic, person.Relative)
	if err != nil {
		return err
	}
	return nil
}
