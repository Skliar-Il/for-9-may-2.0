package service

import (
	"fmt"
	"for9may/internal/dto"
	"for9may/internal/repository"
	"for9may/pkg/database"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"mime/multipart"
	"path/filepath"
)

type GalleryService struct {
	GalleryRepository repository.GalleryRepositoryInterface
	DBPool            *pgxpool.Pool
}

func NewGalleryService(dbPool *pgxpool.Pool, galleryRepository repository.GalleryRepositoryInterface) *GalleryService {
	return &GalleryService{
		GalleryRepository: galleryRepository,
		DBPool:            dbPool,
	}
}

func (g *GalleryService) CreatePost(ctx *gin.Context, post *dto.CreateGalleryPostDTO) (*uuid.UUID, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := g.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, logger.GetLoggerFromCtx(ctx))

	id, err := g.GalleryRepository.Create(ctx, tx, post)
	if err != nil {
		localLogger.Error(ctx, "create gallery post database error", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return id, nil
}

func (g *GalleryService) DeletePost(ctx *gin.Context, postID uuid.UUID) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := g.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, logger.GetLoggerFromCtx(ctx))

	if err := g.GalleryRepository.Delete(ctx, tx, postID); err != nil {
		localLogger.Error(ctx, "delete post database error")
		return web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
		return web.InternalServerError{}
	}

	return nil
}

func (g *GalleryService) GetPosts(ctx *gin.Context) ([]*dto.GalleryPostDTO, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := g.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, logger.GetLoggerFromCtx(ctx))

	posts, err := g.GalleryRepository.GetAll(ctx, tx)
	if err != nil {
		localLogger.Error(ctx, "database error", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return posts, nil
}

func (g *GalleryService) UploadPostFile(ctx *gin.Context, file *multipart.FileHeader, postID uuid.UUID) (string, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := g.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return "", web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, logger.GetLoggerFromCtx(ctx))

	fileName := uuid.New()
	savePath := filepath.Join("upload", fmt.Sprintf("%s.jpg", fileName.String()))
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		localLogger.Error(ctx, "save file error", zap.Error(err))
		return "", web.InternalServerError{}
	}

	databaseFilePath := fmt.Sprintf("/files/%s.jpg", fileName.String())

	if err := g.GalleryRepository.UpdateLink(ctx, tx, postID, databaseFilePath); err != nil {
		localLogger.Error(ctx, "save file error", zap.Error(err))
		return "", nil
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
		return "", web.InternalServerError{}
	}

	return fmt.Sprintf("/files/%s.jpg", fileName.String()), nil
}
