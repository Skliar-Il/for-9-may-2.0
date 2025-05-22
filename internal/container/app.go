package container

import (
	"context"
	"errors"
	"fmt"
	"for9may/internal/config"
	httpserver "for9may/internal/transport/http"
	"for9may/pkg/database"
	jwtservice "for9may/pkg/jwt"
	"for9may/pkg/trace"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"net/http"
	"time"
)

const (
	startServerErrorString    = "failed start server: %v"
	getConfigErrorString      = "failed get config: %v"
	loadPrivateKeyErrorString = "failed load privet key: %v"
	startDataBaseErrorString  = "failed start database: %v"
)

// NewApp
// @title polk sirius
func NewApp() *http.Server {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf(getConfigErrorString, err)
	}
	serverEngine := gin.Default()

	gin.SetMode(cfg.Server.ServerMode)

	serverEngine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{
			"Origin",
			"Accept",
			"X-Requested-With",
			"Content-Type",
			"Authorization",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	server := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", cfg.Server.HttpPort),
		Handler:        serverEngine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("app started on port: %d", cfg.Server.HttpPort)

	privetKey, err := jwtservice.LoadPrivateKey("./certs/private.pem")
	if err != nil {
		log.Fatalf(loadPrivateKeyErrorString, err)
	}
	publicKey, err := jwtservice.LoadPublicKey("./certs/public.pem")
	if err != nil {
		log.Fatalf(loadPrivateKeyErrorString, err)
	}
	serviceJwt := jwtservice.NewServiceJWT(privetKey, publicKey, time.Hour*4, time.Minute*15)

	dbPool, err := database.New(ctx, cfg.DataBase)
	if err != nil {
		log.Fatalf(startDataBaseErrorString, err)
	}

	_, err = trace.InitTracer(fmt.Sprintf("http://%s:%s/api/traces", cfg.Trace.Host, cfg.Trace.Port), "Polk sirius")
	if err != nil {
		log.Fatalf("init tracer error: %v", err)
	}
	serverEngine.Use(otelgin.Middleware("polk-sirius"))

	httpserver.Define(serverEngine, cfg, serviceJwt, dbPool)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(startServerErrorString, err)
		}
	}()

	return server
}
