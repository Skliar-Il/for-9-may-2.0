package http

import (
	"for9may/internal/dto"
	"for9may/internal/service"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// DeleteMedal
// @Tags Medal
// @Param id path int true "medal id"
// Success 204
// Failure 422 web.ValidationError
// Failure 500 web.InternalServerError
// @router /medal/{id} [delete]
func (m *MedalHandler) DeleteMedal(c *gin.Context) {
	medalIDStr := c.Param("id")
	medalID, err := strconv.ParseInt(medalIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "id mast be int")
		return
	}

	if err := m.MedalService.DeleteMedal(c, int(medalID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateMedal
// @Tags Medal
// @Param medal body dto.MedalDTO true "medal"
// Success 204
// Failure 422 web.ValidationError
// Failure 500 web.InternalServerError
// @router /medal [put]
func (m *MedalHandler) UpdateMedal(c *gin.Context) {
	var medal dto.MedalDTO
	if err := c.ShouldBindJSON(&medal); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := m.MedalService.UpdateMedal(c, &medal); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
