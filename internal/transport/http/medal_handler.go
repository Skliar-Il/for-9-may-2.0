package http

import (
	"for9may/internal/dto"
	"for9may/internal/service"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MedalHandler struct {
	MedalService *service.MedalService
}

func NewMedalHandler(medalService *service.MedalService) *MedalHandler {
	return &MedalHandler{MedalService: medalService}
}

// CreateMedal
// @Tags Medal
// @Router /medal/create [post]
// @Param medal body dto.CreateMedalDTO true "Create medal"
// Success 201
// Failure 422
// Failure 500
func (m *MedalHandler) CreateMedal(c *gin.Context) {
	var medal dto.CreateMedalDTO
	if err := c.ShouldBindJSON(&medal); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: err.Error()})
		return
	}

	medalID, err := m.MedalService.CreateMedal(c, &medal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": medalID})
}

// GetMedals
// @Tags Medal
// @router /medal [get]
// Success 200 {array} []dto.MedalDTO
// Failure 422 web.ValidationError
// Failure 500 web.InternalServerError
func (m *MedalHandler) GetMedals(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	medals, err := m.MedalService.GetMedals(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	localLogger.Info(c, "get medals")

	c.JSON(http.StatusOK, medals)
}
