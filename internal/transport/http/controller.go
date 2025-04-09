package http

import (
	_ "for9may/docs"
	"for9may/internal/config"
	"for9may/internal/repository"
	"for9may/internal/service"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Define(engine *gin.Engine, cfg *config.Config, serviceJwt *jwtservice.ServiceJWT, dbPool *pgxpool.Pool) {
	engine.Use(logger.Middleware())
	engine.Use(gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api")

	personRepository := repository.NewPersonRepository()
	medalRepository := repository.NewMedalRepository()

	personService := service.NewPersonService(dbPool, personRepository, medalRepository)

	personController := NewPersonHandler(personService)
	profileController := NewProfileHandler(serviceJwt, cfg.Admin)

	profileGroup := api.Group("/profile")
	{
		profileGroup.POST("/login", profileController.LoginAdmin)
		profileGroup.POST("/refresh", profileController.RefreshAdmin)
	}

	personGroup := api.Group("/person")
	{
		personGroup.POST("/create", personController.NewPerson)
	}
}
