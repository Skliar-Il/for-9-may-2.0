package http

import (
	"fmt"
	"for9may/internal/model"
	"for9may/internal/service"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PersonHandler struct {
	PersonService *service.PersonService
}

func NewPersonHandler(personService *service.PersonService) *PersonHandler {
	return &PersonHandler{PersonService: personService}
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
	}

	c.JSON(201, gin.H{"id": personId})
	return
}
