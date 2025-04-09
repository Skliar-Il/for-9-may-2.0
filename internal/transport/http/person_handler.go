package http

import (
	"errors"
	"fmt"
	"for9may/internal/model"
	"for9may/internal/service"
	"for9may/pkg/jwt"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	PersonService *service.PersonService
	JWTService    *jwt.ServiceJWT
}

func NewPersonHandler(personService *service.PersonService, jwtService *jwt.ServiceJWT) *PersonHandler {
	return &PersonHandler{PersonService: personService, JWTService: jwtService}
}

// NewPerson
// @Tags Person
// @Router /person/create [post]
// @Param person body model.CreatePersonModel true "New Person"
// @Success 201
func (p *PersonHandler) NewPerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	var person model.CreatePersonModel
	if err := c.ShouldBindJSON(&person); err != nil {
		localLogger.Error(c, fmt.Sprintf("invalid body: %v", err))
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

// GetPerson
// @Summary Get person information
// @Description Retrieve person data with optional status check
// @Tags Person
// @Accept json
// @Produce json
// @Param check query boolean false "Status check flag" default(false)
// @Success 200
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /person [get]
func (p *PersonHandler) GetPerson(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	statusStr := c.DefaultQuery("check", "true")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		localLogger.Error(c, "invalid query param", zap.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "field check is invalid"})
	}
	if !status {
		token, err := c.Cookie(jwtservice.AccessTokenCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, web.UnAuthorizedError{})
				return
			}
			localLogger.Error(c, "get token from cookie", zap.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.InternalServerError{})
		}
		claims, err := p.JWTService.DecodeKey(token)
	}
}
