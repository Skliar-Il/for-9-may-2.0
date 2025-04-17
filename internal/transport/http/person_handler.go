package http

import (
	"errors"
	"fmt"
	"for9may/internal/config"
	"for9may/internal/dto"
	"for9may/internal/service"
	"for9may/pkg/jwt"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	PersonService  *service.PersonService
	ProfileService *service.ProfileService
	JWTService     *jwt.ServiceJWT
	PhotoConfig    config.PhotoConfig
}

func NewPersonHandler(
	personService *service.PersonService,
	profileService *service.ProfileService,
	jwtService *jwt.ServiceJWT,
	photoConfig config.PhotoConfig,
) *PersonHandler {
	return &PersonHandler{
		PersonService:  personService,
		JWTService:     jwtService,
		ProfileService: profileService,
		PhotoConfig:    photoConfig,
	}
}

// NewPerson
// @Tags Person
// @Router /person/create [post]
// @Failure 422
// @Param person body dto.CreatePersonDTO true "New Person"
// @Success 201
func (p *PersonHandler) NewPerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	var person dto.CreatePersonDTO
	if err := c.ShouldBindJSON(&person); err != nil {
		localLogger.Info(c, fmt.Sprintf("invalid body: %v", err))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: err.Error()})
		return
	}

	personId, err := p.PersonService.CreatePeron(c, &person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(201, gin.H{"id": personId})
	return
}

// GetPersonList
// @Summary Get person information list
// @Description Retrieve person data with optional status check
// @Tags Person
// @Accept json
// @Produce json
// @Param check query bool true "Status check flag" default(true)
// @Success 200 {array} []dto.PersonDTO
// @Failure 400 {object} web.BadRequestError "Invalid request parameters"
// @Failure 422
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /person [get]
func (p *PersonHandler) GetPersonList(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	statusStr := c.DefaultQuery("check", "true")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		localLogger.Error(c, "invalid query param", zap.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "field check is invalid"})
		return
	}
	if !status {
		token, err := c.Cookie(jwtservice.AccessTokenCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, jwtservice.UndefinedTokenError) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, web.UnAuthorizedError{})
				return
			}
			localLogger.Error(c, "get token from cookie", zap.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.InternalServerError{})
			return
		}
		if err := p.ProfileService.CheckAccount(c, token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
	}

	persons, err := p.PersonService.GetPersons(c, status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, persons)
}

// ValidatePerson
// @Summary ValidatePerson person
// @Description status switch status check on true
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200
// @Failure 400 {object} web.BadRequestError "Invalid request parameters"
// @Failure 422
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /person/validate/{id} [patch]
func (p *PersonHandler) ValidatePerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)
	idStr := c.Param("id")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "invalid path param id"})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		localLogger.Info(c, "parse uuid error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "invalid type id"})
	}
	if err := p.PersonService.ValidatePerson(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
	return
}

// DeletePerson
// @Tags Person
// @Param id path string true "User ID"
// @Success 204
// @Failure 422
// @Failure 500
// @Router /person/{id} [delete]
func (p *PersonHandler) DeletePerson(c *gin.Context) {
	personIDStr := c.Param("id")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			web.ValidationError{Message: "invalid format id, mast be uuid v4"},
		)
		return
	}

	if err := p.PersonService.DeletePerson(c, personID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetPersonByID godoc
// @Summary Get person details
// @Description Retrieves complete information about a person by their ID, including medal awards
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "Person's unique identifier (UUID)" format(uuid)
// @Success 200 {object} dto.PersonDTO "Successfully retrieved person data"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /person/{id} [get]
func (p *PersonHandler) GetPersonByID(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)
	statusAuth := true

	personIDStr := c.Param("id")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "invalid format id, mast be uuid v4")
	}

	token, err := c.Cookie(jwtservice.AccessTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) || errors.Is(err, jwtservice.UndefinedTokenError) {
			statusAuth = false
		} else {
			localLogger.Error(c, "get token from cookie", zap.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, web.InternalServerError{})
			return
		}
	}

	if statusAuth != false {
		if err := p.ProfileService.CheckAccount(c, token); err != nil {
			statusAuth = false
		}
	}

	person, err := p.PersonService.GetPersonByID(c, personID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if person.StatusCheck == false && statusAuth == false {
		c.AbortWithStatusJSON(http.StatusUnauthorized, web.UnAuthorizedError{})
		return
	}

	c.JSON(http.StatusOK, person)
}

// UpdatePerson
// @Summary Update person information
// @Description Updates existing person's data by ID with provided information
// @Tags Person
// @Accept json
// @Produce json
// @Param request body dto.PersonDTO true "Person data to update"
// @Success 204 "No content (successful update with no response body)"
// @Failure 400 {object} web.ValidationError "Invalid request format"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Person not found"
// @Failure 422 {object} web.ValidationError "Validation error"
// @Failure 500 "Internal server error"
// @Router /persons/{id} [put]
func (p *PersonHandler) UpdatePerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	var person dto.PersonDTO
	if err := c.ShouldBindJSON(&person); err != nil {
		localLogger.Error(c, fmt.Sprintf("invalid body: %v", err))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// CountPerson
// @Summery get count not check person
// @Accept json
// @Produce json
// @Tags Person
// @Success 200 {object} dto.PersonCountDTO
// @Failure 401 {object} web.UnAuthorizedError
// @Failure 500 {object} web.InternalServerError
// @Router /person/count [get]
func (p *PersonHandler) CountPerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)
	personCount, err := p.PersonService.CountPerson(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	localLogger.Info(c, "get count unread person")
	c.JSON(http.StatusOK, personCount)
}

// UploadFile
// @Summary upload file
// @Description Upload file to s3 storage, use only .jpg and .png
// @Tags Person
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file (jpeg/png)"
// @Param main query boolean true "use as main photo (default false)"
// @Param id path string true "person id" format(uuid)
// @Success 201 {object} dto.CreateNewPhotoDTO
// @Failure 400
// @Failure 401
// @Failure 413
// @Failure 500
// @Router /person/file/upload/{id} [post]
func (p *PersonHandler) UploadFile(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)
	mainStatusStr := c.DefaultQuery("main", "false")
	personIDStr := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "fail is to load"})
			return
		}

		localLogger.Error(c, "load file error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.InternalServerError{})
		return
	}

	mainStatus, err := strconv.ParseBool(mainStatusStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "status must be boolean"})
		return
	}

	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "id must be uuid v4")
		return
	}

	var personPhotoDTO dto.CreateNewPhotoDTO
	personPhotoDTO.MainStatus = mainStatus
	personPhotoDTO.PersonID = personID
	err = p.PersonService.UploadPersonPhoto(c, &personPhotoDTO, p.PhotoConfig.MaxCount, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, personPhotoDTO)
}
