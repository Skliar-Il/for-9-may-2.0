package http

import (
	"errors"
	"fmt"
	"for9may/internal/config"
	"for9may/internal/models"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	errorGenerateTokenString = "generate token error: %v"
)

type ProfileController struct {
	ServiceJWT *jwtservice.ServiceJWT
	AdminCfg   config.AdminCfg
}

func NewProfileHandler(serviceJwt *jwtservice.ServiceJWT, adminCfg config.AdminCfg) *ProfileController {
	return &ProfileController{
		ServiceJWT: serviceJwt,
		AdminCfg:   adminCfg,
	}
}

// LoginAdmin
// @Summary Login admin
// @Description Authenticate admin with basic auth
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Success 200 {object} models.ProfileLoginResponse "Authorization OK"
// @Failure 401 {object} web.ErrorResponse "Authorization error"
// @Router /profile/login [post]
func (p *ProfileController) LoginAdmin(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		localLogger.Error(c, "invalid format BasicAuth")
		c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.InvalidBasicAuthForm})
		return
	}

	if username == p.AdminCfg.Login && password == p.AdminCfg.Password {
		localLogger.Info(c, "generate tokens")

		refreshToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(username, jwtservice.RefreshToken))
		if err != nil {
			localLogger.Error(c, fmt.Sprintf(errorGenerateTokenString, err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
			return
		}

		accessToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(username, jwtservice.AccessToken))
		if err != nil {
			localLogger.Error(c, fmt.Sprintf(errorGenerateTokenString, err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
			return
		}

		localLogger.Info(c, "push tokens in cookie")

		p.ServiceJWT.SetCookieRefresh(c, refreshToken)
		c.SetSameSite(http.SameSiteStrictMode)
		p.ServiceJWT.SetCookieAccess(c, accessToken)

		c.JSON(http.StatusOK, models.ProfileLoginResponse{Message: accessToken})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.InvalidLoginPassword})
		return
	}
}

// RefreshAdmin
// @Tags auth
// @Summary Refresh admin tokens
// @Description Refresh access and refresh tokens for admin
// @Router /profile/refresh [post]
// @Produce json
// @Success 200 {object} models.ProfileLoginResponse
// @Failure 401 {object} web.ErrorResponse
// @Failure 403 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
func (p *ProfileController) RefreshAdmin(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	token, err := c.Cookie(jwtservice.RefreshTokenCookieName)
	if err != nil {
		if errors.Is(jwtservice.UndefinedTokenError, err) || errors.Is(http.ErrNoCookie, err) {
			localLogger.Error(c, "token is nil")
			c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.TokenExpectedError})
			return
		} else {
			localLogger.Error(c, fmt.Sprintf("get token error: %v", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
			return
		}
	}

	localLogger.Info(c, "decode token")

	tokenClaims, err := p.ServiceJWT.DecodeKey(token)
	if err != nil {
		if errors.Is(jwtservice.InvalidTokenError, err) {
			localLogger.Error(c, fmt.Sprintf("invalid claims: %v", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.TokenInvalidError})
			return
		} else {
			localLogger.Error(c, fmt.Sprintf("failed get claims: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{Message: web.InternalServerError})
			return
		}
	}

	localLogger.Info(c, "validate user")

	if tokenClaims.Subject != p.AdminCfg.Login {
		localLogger.Error(c, "invalid user")
		c.AbortWithStatusJSON(http.StatusForbidden, web.ErrorResponse{Message: web.InvalidSubjectError})
		return
	}

	localLogger.Info(c, "generate tokens")

	refreshToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(tokenClaims.Subject, jwtservice.RefreshToken))
	if err != nil {
		localLogger.Error(c, fmt.Sprintf(errorGenerateTokenString, err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
		return
	}

	accessToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(tokenClaims.Subject, jwtservice.AccessToken))
	if err != nil {
		localLogger.Error(c, fmt.Sprintf(errorGenerateTokenString, err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
		return
	}

	localLogger.Info(c, "push tokens in cookie")

	p.ServiceJWT.SetCookieRefresh(c, refreshToken)
	c.SetSameSite(http.SameSiteStrictMode)
	p.ServiceJWT.SetCookieAccess(c, accessToken)

	c.JSON(http.StatusOK, models.ProfileLoginResponse{Message: accessToken})
}
