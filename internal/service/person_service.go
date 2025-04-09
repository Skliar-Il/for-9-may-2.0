package service

import (
	"context"
	"for9may/internal/model"
	"for9may/internal/repository"
	"for9may/pkg/database"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PersonService struct {
	DBPool           *pgxpool.Pool
	PersonRepository repository.PersonRepositoryInterface
	MedalRepository  repository.MedalRepositoryInterface
}

func NewPersonService(
	dbPool *pgxpool.Pool,
	personRepository repository.PersonRepositoryInterface,
	medalRepository repository.MedalRepositoryInterface,
) *PersonService {
	return &PersonService{
		DBPool:           dbPool,
		PersonRepository: personRepository,
		MedalRepository:  medalRepository,
	}
}

func (p *PersonService) CreatePeron(ctx context.Context, person *model.CreatePersonModel) (*uuid.UUID, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, logger.GetLoggerFromCtx(ctx))

	personUUID, err := p.PersonRepository.CreatePerson(ctx, tx, person)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			localLogger.Error(ctx, "database error: ", zap.String("error", pgError.String()))

			switch pgError.Type {
			case database.TypeDuplicate:
				return nil, web.AlreadyExistError{Message: "person already exist"}

			case database.TypeNotNull:
				return nil, web.BadRequestError{Message: "send null variable"}
			}
			return nil, err
		}
	}
	err = p.MedalRepository.CreateMedalPerson(ctx, tx, personUUID, person.Medals)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			localLogger.Error(ctx, "database error: ", zap.String("error", pgError.String()))
			if pgError.Type == database.TypeForeignKey {
				return nil, web.BadRequestError{Message: "invalid medals id"}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return personUUID, nil
}
