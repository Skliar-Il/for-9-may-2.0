package container

import (
	"errors"
	"fmt"
	"for9may/internal/config"
	httpserver "for9may/internal/transport/http"
	jwtservice "for9may/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	startServerError    = "failed start server: %v"
	getConfigError      = "failed get config: %v"
	loadPrivateKeyError = "failed load privet key: %v"
)

// NewApp
// @title polk sirius
func NewApp() *http.Server {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf(getConfigError, err)
	}
	serverEngine := gin.Default()

	gin.SetMode(cfg.Server.ServerMode)

	serverEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	server := &http.Server{
		Addr:           fmt.Sprintf("localhost:%d", cfg.Server.HttpPort),
		Handler:        serverEngine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	privetKey, err := jwtservice.LoadPrivateKey("./certs/private.pem")
	if err != nil {
		log.Fatalf(loadPrivateKeyError, err)
	}
	publicKey, err := jwtservice.LoadPublicKey("./certs/public.pem")
	if err != nil {
		log.Fatalf(loadPrivateKeyError, err)
	}
	serviceJwt := jwtservice.NewServiceJWT(privetKey, publicKey, time.Hour*4, time.Minute*15)

	httpserver.Define(serverEngine, cfg, serviceJwt)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(startServerError, err)
		}
	}()

	return server
}
