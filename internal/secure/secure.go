package secure

import (
	"errors"
	"for9may/internal/config"
	jwtservice "for9may/pkg/jwt"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"slices"
)

func Middleware(checkLink map[string][]string, jwtService *jwtservice.ServiceJWT, adminCfg config.AdminCfg) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		url := ctx.FullPath()

		checkMethods, exist := checkLink[url]
		if !exist {
			ctx.Next()
			return
		}

		if slices.Contains(checkMethods, method) {
			ctx.Next()
			return
		}

		token, err := ctx.Cookie(jwtservice.AccessTokenCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, jwtservice.UndefinedTokenError) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, web.ForbiddenError{})
				return
			} else {
				log.Printf("get token claims error: %v", err)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.InternalServerError{})
				return
			}
		}

		claims, err := jwtService.DecodeKey(token)
		if err != nil {
			if errors.Is(jwtservice.InvalidTokenError, err) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{Message: web.TokenInvalidErrorString})
				return
			} else {
				return
			}
		}
		if claims.Subject != adminCfg.Login {
			return
		}

		ctx.Next()
	}
}
