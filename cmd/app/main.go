package main

import (
	"context"
	"fmt"
	"for9may/internal/container"
	"os"
	"os/signal"
	"syscall"
)

// @title Polk Sirius
// @securityDefinitions.basic  BasicAuth
// @version 2.0
// @BasePath /api
func main() {
	ctx := context.Background()
	stopChanel := make(chan os.Signal, 1)

	server := container.NewApp()

	signal.Notify(stopChanel, syscall.SIGTERM, syscall.SIGINT)
	<-stopChanel

	server.Shutdown(ctx)

	fmt.Println("Server Stopped")
}
