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

func Define(engine *gin.Engine, cfg *config.Config, jwtService *jwtservice.ServiceJWT, dbPool *pgxpool.Pool) {
	mainLogger := logger.New()

	engine.Use(logger.Middleware(mainLogger))
	engine.Use(gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api")

	personRepository := repository.NewPersonRepository()
	medalRepository := repository.NewMedalRepository()
	formRepository := repository.NewFormRepository()
	ownerRepository := repository.NewOwnerRepository()

	personService := service.NewPersonService(
		dbPool,
		personRepository,
		medalRepository,
		formRepository,
		ownerRepository,
	)
	profileService := service.NewProfileService(cfg.Admin, jwtService)
	medalService := service.NewMedalService(dbPool, medalRepository)

	personController := NewPersonHandler(personService, profileService, jwtService)
	profileController := NewProfileHandler(jwtService, cfg.Admin)
	medalController := NewMedalHandler(medalService)

	profileGroup := api.Group("/profile")
	{
		profileGroup.POST("/login", profileController.LoginAdmin)
		profileGroup.POST("/refresh", profileController.RefreshAdmin)
	}

	personGroup := api.Group("/person")
	{
		personGroup.POST("/create", personController.NewPerson)
		personGroup.GET("", personController.GetPersonList)
		personGroup.PATCH("/validate/:id", personController.ValidatePerson)
		personGroup.DELETE("/:id", personController.DeletePerson)
		personGroup.GET("/:id", personController.GetPersonByID)
		personGroup.PUT("", personController.UpdatePerson)
		personGroup.GET("/count", personController.CountPerson)
	}

	medalGroup := api.Group("/medal")
	{
		medalGroup.POST("/create", medalController.CreateMedal)
		medalGroup.GET("", medalController.GetMedals)
	}
}
