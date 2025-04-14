package service

import (
	"fmt"
	"for9may/internal/model"
	"for9may/internal/repository"
	"for9may/pkg/database"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MedalService struct {
	MedalRepository repository.MedalRepositoryInterface
	DBPool          *pgxpool.Pool
}

func NewMedalService(dbPoll *pgxpool.Pool, medalRepository repository.MedalRepositoryInterface) *MedalService {
	return &MedalService{
		MedalRepository: medalRepository,
		DBPool:          dbPoll,
	}
}

func (m *MedalService) GetMedals(ctx *gin.Context) ([]model.MedalModel, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := m.DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx error: %w", err)
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	medals, err := m.MedalRepository.GetMedals(ctx, tx)
	if err != nil {
		localLogger.Error(ctx, "get medals error", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return medals, nil
}

func (m *MedalService) CreateMedal(ctx *gin.Context, medal *model.CreateMedalModel) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := m.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx error: %w", err)
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	if err := m.MedalRepository.CreateMedal(ctx, tx, medal); err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			if pgError.Type == database.TypeDuplicate {
				return web.BadRequestError{Message: "medal already exist"}

			} else {

				localLogger.Error(ctx, "database error", zap.Error(err))
				return web.InternalServerError{}
			}
		}

		localLogger.Error(ctx, "repository error", zap.Error(err))
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
	}

	return nil
}
