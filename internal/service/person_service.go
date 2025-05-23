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

type PersonService struct {
	DBPool           *pgxpool.Pool
	PersonRepository repository.PersonRepositoryInterface
	MedalRepository  repository.MedalRepositoryInterface
	FormRepository   repository.FormRepositoryInterface
	OwnerRepository  repository.OwnerRepositoryInterface
	PhotoRepository  repository.PhotoRepositoryInterface
}

func NewPersonService(
	dbPool *pgxpool.Pool,
	personRepository repository.PersonRepositoryInterface,
	medalRepository repository.MedalRepositoryInterface,
	formRepository repository.FormRepositoryInterface,
	ownerRepository repository.OwnerRepositoryInterface,
	photoRepository repository.PhotoRepositoryInterface,
) *PersonService {
	return &PersonService{
		DBPool:           dbPool,
		PersonRepository: personRepository,
		MedalRepository:  medalRepository,
		FormRepository:   formRepository,
		OwnerRepository:  ownerRepository,
		PhotoRepository:  photoRepository,
	}
}

func (p *PersonService) CreatePeron(ctx *gin.Context, person *dto.CreatePersonDTO) (*uuid.UUID, error) {
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
	err = p.MedalRepository.CreateMedalPerson(ctx, tx, *personUUID, person.Medals)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			localLogger.Error(
				ctx,
				"create relative medal person database error: ",
				zap.String("error", pgError.String()),
			)
			if pgError.Type == database.TypeForeignKey {
				return nil, web.BadRequestError{Message: "invalid medals id"}
			}
		}
	}

	formID, err := p.FormRepository.CreateForm(ctx, tx, personUUID)
	if err != nil {
		localLogger.Error(ctx, "create form database error", zap.String("error", err.Error()))
		return nil, web.InternalServerError{}
	}

	if err := p.OwnerRepository.Create(ctx, tx, person, formID); err != nil {
		localLogger.Error(ctx, "create owner database error: ", zap.String("error", err.Error()))

	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return personUUID, nil
}

func (p *PersonService) GetPersons(ctx *gin.Context, check bool, status bool) ([]*dto.PersonDTO, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error: ", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	persons, err := p.PersonRepository.GetPersons(ctx, tx, check, status)
	if err != nil {
		localLogger.Error(ctx, "get person error", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	if persons == nil {
		persons = []*dto.PersonDTO{}
	}

	return persons, nil
}

func (p *PersonService) CountPerson(ctx *gin.Context) (*dto.PersonCountDTO, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	personCount, err := p.PersonRepository.CountUnread(ctx, tx)
	if err != nil {
		localLogger.Error(ctx, "get count unread person error", zap.Error(err))
		return nil, web.InternalServerError{}
	}

	return personCount, nil
}

func (p *PersonService) GetPersonByID(ctx *gin.Context, personID uuid.UUID) (*dto.PersonDTO, error) {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "start tx error", zap.Error(err))
		return nil, web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	person, err := p.PersonRepository.GerPersonByID(ctx, tx, personID)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			if pgError.Type == database.TypeNoRows {
				return nil, web.NotFoundError{Message: "person not found"}
			}
		}
		localLogger.Error(ctx, "get person error", zap.Error(err))
	}

	return person, nil
}

func (p *PersonService) ValidatePerson(ctx *gin.Context, id uuid.UUID) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "start tx error", zap.Error(err))
		return web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	err = p.PersonRepository.Validate(ctx, tx, id)
	if err != nil {
		localLogger.Error(ctx, "database error", zap.Error(err))
		return web.InternalServerError{}
	}
	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error: ", zap.Error(err))
		return web.InternalServerError{}
	}

	return nil
}

func (p *PersonService) DeletePerson(ctx *gin.Context, id uuid.UUID) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	if err := p.PersonRepository.Delete(ctx, tx, id); err != nil {
		if pgError := database.ValidatePgxError(err); pgError != nil {
			if pgError.Type == database.TypeNoRows {
				return web.NotFoundError{Message: "medal not found"}
			}
		}

		localLogger.Error(ctx, "database error", zap.Error(err))
		return web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
	}

	return nil
}

func (p *PersonService) UploadPersonPhoto(
	ctx *gin.Context,
	photo *dto.CreatePhotoDTO,
	countOK int,
	file *multipart.FileHeader,
) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
		return web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	status, err := p.PhotoRepository.CheckCount(ctx, tx, countOK, photo.PersonID)
	if err != nil {
		localLogger.Error(ctx, "check count photo error", zap.Error(err))
		return web.InternalServerError{}
	}

	if status == false {
		return web.BadRequestError{Message: fmt.Sprintf("to many photo, max: %d", countOK)}
	}

	if photo.MainStatus == true {
		mainStatusExist, err := p.PhotoRepository.CheckMainStatus(ctx, tx, photo.PersonID)
		if err != nil {
			localLogger.Error(ctx, "check main status error", zap.Error(err))
			return web.InternalServerError{}
		}
		if mainStatusExist == true {
			return web.BadRequestError{Message: "main photo already exist"}
		}
	}

	fileName := uuid.New()
	link := fmt.Sprintf("/files/%s.jpg", fileName.String())
	photo.Link = link
	savePath := filepath.Join("upload", fmt.Sprintf("%s.jpg", fileName.String()))
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		localLogger.Error(ctx, "save file error", zap.Error(err))
		return web.InternalServerError{}
	}

	err = p.PhotoRepository.CreatePhoto(ctx, tx, photo)
	if err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			if pgError.Type == database.TypeForeignKey {
				return web.NotFoundError{Message: "person not found"}
			}
		}
		localLogger.Error(ctx, "failed add photo data in database", zap.Error(err))
		return web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
		return web.InternalServerError{}
	}

	return nil
}

func (p *PersonService) DeletePersonPhoto(ctx *gin.Context, photoID int) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
		return web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	if err := p.PhotoRepository.DeletePhoto(ctx, tx, photoID); err != nil {
		localLogger.Error(ctx, "database delete person photo error", zap.Error(err))
		return web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit error", zap.Error(err))
		return web.InternalServerError{}
	}

	return nil
}

func (p *PersonService) UpdatePerson(ctx *gin.Context, person *dto.UpdatePersonDTO) error {
	localLogger := logger.GetLoggerFromCtx(ctx)
	tx, err := p.DBPool.Begin(ctx)
	if err != nil {
		localLogger.Error(ctx, "begin tx error", zap.Error(err))
		return web.InternalServerError{}
	}
	defer database.RollbackTx(ctx, tx, localLogger)

	if err := p.PersonRepository.Update(ctx, tx, person); err != nil {
		localLogger.Error(ctx, "database update person error", zap.Error(err))
		return web.InternalServerError{}
	}
	if err := p.MedalRepository.DeleteMedalPersonByID(ctx, tx, person.ID); err != nil {
		localLogger.Error(ctx, "database delete medal person error", zap.Error(err))
		return web.InternalServerError{}
	}
	if err := p.MedalRepository.CreateMedalPerson(ctx, tx, person.ID, person.Medals); err != nil {
		pgError := database.ValidatePgxError(err)
		if pgError != nil {
			if pgError.Type == database.TypeForeignKey {
				return web.BadRequestError{Message: "invalid medal id"}
			}
		}
		localLogger.Error(ctx, "database create medal person error", zap.Error(err))
		return web.InternalServerError{}
	}
	if err := p.OwnerRepository.Update(ctx, tx, person); err != nil {
		localLogger.Error(ctx, "update owner error", zap.Error(err))
		return web.InternalServerError{}
	}
	if err := p.FormRepository.UpdateForm(ctx, tx, person.ID, person.StatusMain); err != nil {
		localLogger.Error(ctx, "database update form error", zap.Error(err))
		return web.InternalServerError{}
	}

	if err := tx.Commit(ctx); err != nil {
		localLogger.Error(ctx, "commit update person error", zap.Error(err))
	}

	return nil
}
