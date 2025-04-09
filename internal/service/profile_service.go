package service

import (
	"errors"
	"fmt"
	"for9may/internal/config"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileService struct {
	JWTService *jwtservice.ServiceJWT
	AdminCfg   config.AdminCfg
}

func NewProfileService(adminCfg config.AdminCfg, jwtService *jwtservice.ServiceJWT) *ProfileService {
	return &ProfileService{
		JWTService: jwtService,
		AdminCfg:   adminCfg,
	}
}

func (p *ProfileService) CheckAccount(c *gin.Context, token string) error {
	localLogger := logger.GetLoggerFromCtx(c)

	claims, err := p.JWTService.DecodeKey(token)
	if err != nil {
		if errors.Is(jwtservice.InvalidTokenError, err) {
			localLogger.Error(c, fmt.Sprintf("invalid claims: %v", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.TokenInvalidErrorString})
			return web.UnAuthorizedError{}
		} else {
			localLogger.Error(c, fmt.Sprintf("failed get claims: %v", err))
			return web.InternalServerError{}
		}
	}
	if claims.Subject != p.AdminCfg.Login {
		localLogger.Error(c, "invalid login")
		return web.UnAuthorizedError{}
	}
	return nil
}
