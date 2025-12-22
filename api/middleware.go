package api

import (
	"errors"
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/gin-gonic/gin"
	"github.com/orgs/PRPO-skupina-02/Spored/models"
)

const (
	contextTheaterKey = "theater"
)

func SetContextTheater(c *gin.Context, theater models.Theater) {
	c.Set(contextTheaterKey, theater)
}

func GetContextTheater(c *gin.Context) models.Theater {
	tx, ok := c.Get(contextTheaterKey)
	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Could not get theater from context"))
		return models.Theater{}
	}

	return tx.(models.Theater)
}

func TheaterContextMiddleware(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	id, err := request.GetUUIDParam(c, "theaterID")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	theater, err := models.GetTheater(tx, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	SetContextTheater(c, theater)

	c.Next()
}

func TheaterPermissionsMiddleware(c *gin.Context) {
	c.Next()
}
