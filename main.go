package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/orgs/PRPO-skupina-02/Spored/api"
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

	router := gin.Default()
	api.Register(router)

	err := router.Run(":8080")
	if err != nil {
		return err
	}

	return nil
}
