package http

import (
	_ "for9may/docs"
	"for9may/internal/config"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Define(engine *gin.Engine, cfg *config.Config, serviceJwt *jwtservice.ServiceJWT) {
	engine.Use(logger.Middleware())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api")

	profileController := NewProfileHandler(serviceJwt, cfg.Admin)
	profileGroup := api.Group("/profile")
	{
		profileGroup.POST("/login", profileController.LoginAdmin)
		profileGroup.POST("/refresh", profileController.RefreshAdmin)
	}
}
