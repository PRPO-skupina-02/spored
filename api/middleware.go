package api

import (
	"errors"
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/PRPO-skupina-02/spored/models"
	"github.com/gin-gonic/gin"
)

const (
	contextTheaterKey = "theater"
	contextMovieKey   = "movie"
)

func SetContextTheater(c *gin.Context, theater models.Theater) {
	c.Set(contextTheaterKey, theater)
}

func GetContextTheater(c *gin.Context) models.Theater {
	theater, ok := c.Get(contextTheaterKey)
	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Could not get theater from context"))
		return models.Theater{}
	}

	return theater.(models.Theater)
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

func SetContextMovie(c *gin.Context, movie models.Movie) {
	c.Set(contextMovieKey, movie)
}

func GetContextMovie(c *gin.Context) models.Movie {
	movie, ok := c.Get(contextMovieKey)
	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Could not get movie from context"))
		return models.Movie{}
	}

	return movie.(models.Movie)
}

func MovieContextMiddleware(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	id, err := request.GetUUIDParam(c, "movieID")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	movie, err := models.GetMovie(tx, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	SetContextMovie(c, movie)

	c.Next()
}
