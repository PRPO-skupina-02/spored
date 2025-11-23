package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {

	// Healthcheck
	router.GET("/healthcheck", healthcheck)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
