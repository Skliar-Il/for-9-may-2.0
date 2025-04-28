package http

import (
	_ "for9may/docs"
	"for9may/internal/config"
	"for9may/internal/repository"
	"for9may/internal/secure"
	"for9may/internal/service"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var checkLink = map[string][]string{
	"/api/person/validate/:id":    {"PATCH"},
	"/api/person/:id":             {"DELETE"},
	"/api/person":                 {"PUT"},
	"/api/person/count":           {"GET"},
	"/api/person/file/delete/:id": {"DELETE"},

	"/api/medal/create": {"POST"},
	"/api/medal":        {"PUT"},
	"/api/medal/:id":    {"DELETE"},

	"/api/gallery":             {"POST"},
	"/api/gallery/file/upload": {"POST"},
	"/api/gallery/:id":         {"DELETE"},
}

func Define(engine *gin.Engine, cfg *config.Config, jwtService *jwtservice.ServiceJWT, dbPool *pgxpool.Pool) {
	mainLogger := logger.New()

	engine.MaxMultipartMemory = 10 << 20
	engine.Use(logger.Middleware(mainLogger))
	engine.Use(gin.Recovery())
	engine.Use(secure.Middleware(checkLink, jwtService, cfg.Admin))
	engine.Static("/files", "./upload")

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api")

	personRepository := repository.NewPersonRepository()
	medalRepository := repository.NewMedalRepository()
	formRepository := repository.NewFormRepository()
	ownerRepository := repository.NewOwnerRepository()
	photoRepository := repository.NewPhotoRepository()
	galleryRepository := repository.NewGalleryRepository()

	personService := service.NewPersonService(
		dbPool,
		personRepository,
		medalRepository,
		formRepository,
		ownerRepository,
		photoRepository,
	)
	profileService := service.NewProfileService(cfg.Admin, jwtService)
	medalService := service.NewMedalService(dbPool, medalRepository)
	galleryService := service.NewGalleryService(dbPool, galleryRepository)

	personController := NewPersonHandler(personService, profileService, jwtService, cfg.PhotoConfig)
	profileController := NewProfileHandler(jwtService, cfg.Admin)
	medalController := NewMedalHandler(medalService)
	galleryController := NewGalleryHandler(galleryService)

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
		personGroup.POST("/file/upload/:id", personController.UploadFile)
		personGroup.DELETE("/file/delete/:id", personController.DeleteFile)
	}

	medalGroup := api.Group("/medal")
	{
		medalGroup.POST("/create", medalController.CreateMedal)
		medalGroup.GET("", medalController.GetMedals)
		medalGroup.DELETE("/:id", medalController.DeleteMedal)
		medalGroup.PUT("", medalController.UpdateMedal)
	}

	galleryGroup := api.Group("/gallery")
	{
		galleryGroup.POST("", galleryController.CreatePost)
		galleryGroup.DELETE("/:id", galleryController.DeletePost)
		galleryGroup.GET("", galleryController.GetPosts)
		galleryGroup.POST("/file/upload/:id", galleryController.UploadPostFile)
	}
}
