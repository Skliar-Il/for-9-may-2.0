package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"for9may/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PersonRepositoryInterface interface {
	CreatePerson(ctx context.Context, tx pgx.Tx, person *dto.CreatePersonDTO) (*uuid.UUID, error)
	GetPersons(ctx context.Context, tx pgx.Tx, check bool) ([]dto.PersonDTO, error)
	Validate(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	CountUnread(ctx context.Context, tx pgx.Tx) (*dto.PersonCountDTO, error)
	GerPersonByID(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (*dto.PersonDTO, error)
	Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	Update(ctx context.Context, tx pgx.Tx, person *dto.UpdatePersonDTO) error
}

type PersonRepository struct{}

func NewPersonRepository() *PersonRepository {
	return &PersonRepository{}
}

func (PersonRepository) CreatePerson(
	ctx context.Context,
	tx pgx.Tx,
	person *dto.CreatePersonDTO,
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

func (PersonRepository) GetPersons(ctx context.Context, tx pgx.Tx, check bool) ([]dto.PersonDTO, error) {
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

	var persons []dto.PersonDTO
	if rows != nil {
		for rows.Next() {
			var p dto.PersonDTO
			var medalsJSON []byte
			var photoJSON []byte

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
				&p.DatePublished,
				&medalsJSON,
				&photoJSON,
			)
			if err != nil {
				return nil, fmt.Errorf("scan failed: %w", err)
			}

			if err := json.Unmarshal(medalsJSON, &p.Medals); err != nil {
				return nil, fmt.Errorf("failed to unmarshal medals: %w", err)
			}

			if err := json.Unmarshal(photoJSON, &p.Photo); err != nil {
				return nil, fmt.Errorf("failed to unmarshal photo: %w", err)
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

func (PersonRepository) GerPersonByID(ctx context.Context, tx pgx.Tx, personID uuid.UUID) (*dto.PersonDTO, error) {
	query := `
	SELECT *
	FROM all_person_fields_view
	WHERE id = $1
	`
	var p dto.PersonDTO
	var medalsJSON []byte
	var photoJSON []byte
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
		&p.DatePublished,
		&medalsJSON,
		&photoJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	if err := json.Unmarshal(medalsJSON, &p.Medals); err != nil {
		return nil, fmt.Errorf("failed to unmarshal medals: %w", err)
	}

	if err := json.Unmarshal(photoJSON, &p.Photo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal photo: %w", err)
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
			DELETE FROM person p
			WHERE p.id = $1
`
	status, err := tx.Exec(ctx, query, id)
	if err != nil {
		return nil
	}
	if status.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (PersonRepository) CountUnread(ctx context.Context, tx pgx.Tx) (*dto.PersonCountDTO, error) {
	query := `
		SELECT COUNT(*) AS count
		FROM person p
		LEFT JOIN form f ON f.person_id = p.id
		WHERE f.status_check = false
		`
	var personCount dto.PersonCountDTO
	if err := tx.QueryRow(ctx, query).Scan(&personCount.Count); err != nil {
		return nil, err
	}
	return &personCount, nil
}

func (PersonRepository) Update(ctx context.Context, tx pgx.Tx, person *dto.UpdatePersonDTO) error {
	query := `
		UPDATE person 
		SET
			name = $1,
			surname = $2,
			patronymic = $3,
			date_death = $4,
			date_birth = $5,
			city_birth = $6,
			history = $7,
			rank = $8
		WHERE id = $9
	`

	_, err := tx.Exec(ctx, query, person.Name, person.Surname, person.Patronymic, person.DateDeath,
		person.DateBirth, person.City, person.History, person.Rank, person.ID)
	if err != nil {
		return err
	}
	return nil
}
