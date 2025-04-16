package service

import (
	"fmt"
	"for9may/internal/dto"
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

func (m *MedalService) GetMedals(ctx *gin.Context) ([]dto.MedalDTO, error) {
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
	if medals == nil {
		medals = []dto.MedalDTO{}
	}

	return medals, nil
}

func (m *MedalService) CreateMedal(ctx *gin.Context, medal *dto.CreateMedalDTO) (int, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := m.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
		return 0, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	medalID, err := m.MedalRepository.CreateMedal(ctx, tx, medal)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			if pgError.Type == database.TypeDuplicate {
				return 0, web.BadRequestError{Message: "medal already exist"}

			} else {

				localLogger.Error(ctx, "database error", zap.Error(err))
				return 0, web.InternalServerError{}
			}
		}

		localLogger.Error(ctx, "repository error", zap.Error(err))
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
	}

	return medalID, nil
}
