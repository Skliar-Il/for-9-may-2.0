package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"for9may/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PersonRepositoryInterface interface {
	CreatePerson(ctx context.Context, tx pgx.Tx, person *model.CreatePersonModel) (*uuid.UUID, error)
	GetPerson(ctx context.Context, tx pgx.Tx, check bool) ([]model.PersonModel, error)
	Validate(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	CountUnread(ctx context.Context, tx pgx.Tx) (*model.PersonCountModel, error)
	GerPersonByID(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (*model.PersonModel, error)
}

type PersonRepository struct{}

func NewPersonRepository() *PersonRepository {
	return &PersonRepository{}
}

func (PersonRepository) CreatePerson(
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

func (PersonRepository) GetPerson(ctx context.Context, tx pgx.Tx, check bool) ([]model.PersonModel, error) {
	query := `
		SELECT *
		FROM all_person_fields_view
        WHERE status_check = $1
    `

	rows, err := tx.Query(ctx, query, check)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var persons []model.PersonModel
	if rows != nil {
		for rows.Next() {
			var p model.PersonModel
			var medalsJSON []byte

			err := rows.Scan(
				&p.ID,
				&p.Name,
				&p.Surname,
				&p.Patronymic,
				&p.DateBirth,
				&p.DateDeath,
				&p.City,
				&p.History,
				&p.Rank,
				&p.ContactEmail,
				&p.ContactName,
				&p.ContactSurname,
				&p.ContactPatronymic,
				&p.ContactTelegram,
				&p.Relative,
				&p.StatusCheck,
				&medalsJSON,
			)
			if err != nil {
				return nil, fmt.Errorf("scan failed: %w", err)
			}

			if err := json.Unmarshal(medalsJSON, &p.Medals); err != nil {
				return nil, fmt.Errorf("failed to unmarshal medals: %w", err)
			}

			persons = append(persons, p)
		}

	} else {
		return nil, nil
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return persons, nil
}

func (PersonRepository) GerPersonByID(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (*model.PersonModel, error) {
	query := `
	SELECT *
	FROM all_person_fields_view
	WHERE id = $1
	`
	var p model.PersonModel
	var medalsJSON []byte
	err := tx.QueryRow(ctx, query, personID).Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Patronymic,
		&p.DateBirth,
		&p.DateDeath,
		&p.City,
		&p.History,
		&p.Rank,
		&p.ContactEmail,
		&p.ContactName,
		&p.ContactSurname,
		&p.ContactPatronymic,
		&p.ContactTelegram,
		&p.Relative,
		&p.StatusCheck,
		&medalsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}
	if err := json.Unmarshal(medalsJSON, &p.Medals); err != nil {
		return nil, fmt.Errorf("failed to unmarshal medals: %w", err)
	}
	return &p, nil
}

func (PersonRepository) Validate(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	query := `
			UPDATE form f
			SET status_check = true
			FROM person p
			WHERE p.id = f.person_id AND p.id = $1
`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (PersonRepository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	query := `
			DELETE person p
			WHERE p.id = $1
`
	status, err := tx.Exec(ctx, query, id)
	if err != nil {
		return nil
	}
	if status.RowsAffected() != 1 {

	}
	return nil
}

func (PersonRepository) CountUnread(ctx context.Context, tx pgx.Tx) (*model.PersonCountModel, error) {
	query := `
		SELECT COUNT(*) AS count
		FROM person p
		LEFT JOIN form f ON f.person_id = p.id
		WHERE f.status_check = false
		`
	var personCount model.PersonCountModel
	if err := tx.QueryRow(ctx, query).Scan(&personCount.Count); err != nil {
		return nil, err
	}
	return &personCount, nil
}
