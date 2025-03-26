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
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const (
	generateTokenError = "generate token error: %v"
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
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		logger.GetLoggerFromCtx(c).Error(c, "invalid format BasicAuth")
		c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.UnauthorizedError})
	}
	if username == p.AdminCfg.Login && password == p.AdminCfg.Password {
		refreshToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(username, jwtservice.RefreshToken))
		if err != nil {
			logger.GetLoggerFromCtx(c).Error(c, fmt.Sprintf(generateTokenError, err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
		}
		accessToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(username, jwtservice.AccessToken))
		if err != nil {
			logger.GetLoggerFromCtx(c).Error(c, fmt.Sprintf(generateTokenError, err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
		}

		p.ServiceJWT.SetCookieRefresh(c, refreshToken)
		c.SetSameSite(http.SameSiteStrictMode)

		p.ServiceJWT.SetCookieAccess(c, accessToken)

		c.JSON(http.StatusOK, models.ProfileLoginResponse{Message: accessToken})
	} else {
		c.JSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.UnauthorizedError})
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
	token, err := c.Cookie(jwtservice.RefreshTokenCookieName)
	if err != nil {
		if errors.Is(jwtservice.UndefinedTokenError, err) || errors.Is(http.ErrNoCookie, err) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.UnauthorizedError})
			return
		} else {
			logger.GetLoggerFromCtx(c).Error(c, fmt.Sprintf("get token error: %v", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
			return
		}
	}

	tokenClaims, err := p.ServiceJWT.DecodeKey(token)
	if errors.Is(jwt.ErrInvalidKey, err) {
		fmt.Printf("invalid key")
	} else {
		fmt.Printf("penis")
	}

	if tokenClaims.Subject != p.AdminCfg.Login {
		c.AbortWithStatusJSON(http.StatusForbidden, web.ErrorResponse{Message: web.ForbiddenError})
	}

	refreshToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(tokenClaims.Subject, jwtservice.RefreshToken))
	if err != nil {
		logger.GetLoggerFromCtx(c).Error(c, fmt.Sprintf(generateTokenError, err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
	}

	accessToken, err := p.ServiceJWT.Encode(p.ServiceJWT.GetClaims(tokenClaims.Subject, jwtservice.AccessToken))
	if err != nil {
		logger.GetLoggerFromCtx(c).Error(c, fmt.Sprintf(generateTokenError, err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.ErrorResponse{Message: web.InternalServerError})
	}

	p.ServiceJWT.SetCookieRefresh(c, refreshToken)
	c.SetSameSite(http.SameSiteStrictMode)
	p.ServiceJWT.SetCookieAccess(c, accessToken)

	c.JSON(http.StatusOK, models.ProfileLoginResponse{Message: accessToken})
}
