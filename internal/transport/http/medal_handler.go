package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MedalHandler struct {
}

func NewMedalHandler() *MedalHandler {
	return &MedalHandler{}
}

// CreateMedal
// @Tags Medal
// @Router /medal/create [post]
// @Param medal body model.CreateMedalModel true "Create medal"
// Success 201
// Failure 422
// Failure 5oo
func (m *MedalHandler) CreateMedal(c *gin.Context) {
	c.JSON(http.StatusCreated, "пж создай через бд, сегодня доделаю")
}
