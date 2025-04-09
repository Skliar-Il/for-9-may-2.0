package repository

import (
	"context"
	"for9may/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PersonRepositoryInterface interface {
	CreatePerson(ctx context.Context, tx pgx.Tx, person *model.CreatePersonModel) (*uuid.UUID, error)
}

type PersonRepository struct{}

func NewPersonRepository() *PersonRepository {
	return &PersonRepository{}
}

func (p *PersonRepository) CreatePerson(
	ctx context.Context,
	tx pgx.Tx,
	person *model.CreatePersonModel,
) (*uuid.UUID, error) {
	query := `
	INSERT INTO person (
		name, surname, patronymic, date_death, date_birth, city_birth, history, rank
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	) 
	RETURNING id`

	var personID uuid.UUID

	err := tx.QueryRow(ctx, query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.DateDeath,
		person.DateBirth,
		person.City,
		person.History,
		person.Rank,
	).Scan(&personID)
	if err != nil {
		return nil, err
	}

	return &personID, nil
}
