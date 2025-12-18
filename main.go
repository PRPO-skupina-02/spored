package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/PRPO-skupina-02/common/config"
	"github.com/PRPO-skupina-02/common/database"
	"github.com/gin-gonic/gin"
	"github.com/orgs/PRPO-skupina-02/Spored/api"
	"github.com/orgs/PRPO-skupina-02/Spored/db"
)

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("Starting server")

	var logger *slog.Logger

	logLevelConfig := config.GetEnvDefault("LOG_LEVEL", "INFO")
	var logLevel = new(slog.LevelVar)
	switch logLevelConfig {
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}
	slog.Info(fmt.Sprintf("Log level: %s", logLevel.Level()))
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger = slog.New(logHandler)
	slog.SetDefault(logger)

	db, err := database.OpenAndMigrateProd(db.MigrationsFS)
	if err != nil {
		return err
	}

	trans, err := api.RegisterValidation()
	if err != nil {
		return err
	}

	router := gin.Default()
	api.Register(router, db, trans)

	slog.Info("Server startup complete")
	err = router.Run(":8080")
	if err != nil {
		return err
	}

	return nil
}
