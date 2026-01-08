package api

import (
	"net/http"
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/PRPO-skupina-02/spored/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MovieResponse struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"name"`
	Description   string    `json:"description"`
	ImageURL      string    `json:"image_url"`
	Rating        float64   `json:"rating"`
	LengthMinutes int       `json:"length_minutes"`
	Active        bool      `json:"active"`
}

func newMovieResponse(movie models.Movie) MovieResponse {
	return MovieResponse{
		ID:            movie.ID,
		CreatedAt:     movie.CreatedAt,
		UpdatedAt:     movie.UpdatedAt,
		Title:         movie.Title,
		Description:   movie.Description,
		ImageURL:      movie.ImageURL,
		Rating:        movie.Rating,
		LengthMinutes: movie.LengthMinutes,
		Active:        movie.Active,
	}
}

// MoviesList
//
//	@Id				MoviesList
//	@Summary		List movies
//	@Description	List movies
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit the number of responses"	Default(10)
//	@Param			offset	query		int		false	"Offset the first response"		Default(0)
//	@Param			sort	query		string	false	"Sort results"
//	@Success		200		{object}	[]MovieResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/movies [get]
func MoviesList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	pagination := request.GetNormalizedPaginationArgs(c)
	sort := request.GetSortOptions(c)

	movies, total, err := models.GetMovies(tx, pagination, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []MovieResponse{}

	for _, movie := range movies {
		response = append(response, newMovieResponse(movie))
	}

	request.RenderPaginatedResponse(c, response, total)
}

type MovieRequest struct {
	Title         string  `json:"title" binding:"required,min=3"`
	Description   string  `json:"description" binding:"required,min=10"`
	ImageURL      string  `json:"image_url" binding:"required,url"`
	Rating        float64 `json:"rating" binding:"required,min=0,max=10"`
	LengthMinutes int     `json:"length_minutes" binding:"required,min=10,max=1000"`
	Active        bool    `json:"active" binding:"boolean"`
}

// MoviesCreate
//
//	@Id				MoviesCreate
//	@Summary		Create movie
//	@Description	Create movie
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			request	body		MovieRequest	true	"request body"
//	@Success		200		{object}	MovieResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/movies [post]
func MoviesCreate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)

	var req MovieRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	movie := models.Movie{
		ID:            uuid.New(),
		Title:         req.Title,
		Description:   req.Description,
		ImageURL:      req.ImageURL,
		Rating:        req.Rating,
		LengthMinutes: req.LengthMinutes,
		Active:        req.Active,
	}

	err = movie.Create(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, newMovieResponse(movie))
}

// MoviesShow
//
//	@Id				MoviesShow
//	@Summary		Show movie
//	@Description	Show movie
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			movieID	path		string	true	"Movie ID"	Format(uuid)
//	@Success		200		{object}	MovieResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/movies/{movieID} [get]
func MoviesShow(c *gin.Context) {
	movie := GetContextMovie(c)
	c.JSON(http.StatusOK, newMovieResponse(movie))
}

// MoviesUpdate
//
//	@Id				MoviesUpdate
//	@Summary		Update movie
//	@Description	Update movie
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			movieID	path		string			true	"Movie ID"	Format(uuid)
//	@Param			request	body		MovieRequest	true	"request body"
//	@Success		200		{object}	MovieResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/movies/{movieID} [put]
func MoviesUpdate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	movie := GetContextMovie(c)

	var req MovieRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.ImageURL = req.ImageURL
	movie.Rating = req.Rating
	movie.LengthMinutes = req.LengthMinutes
	movie.Active = req.Active

	err = movie.Save(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newMovieResponse(movie))
}

// MoviesDelete
//
//	@Id				MoviesDelete
//	@Summary		Delete movie
//	@Description	Delete movie
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			movieID	path	string	true	"Movie ID"	Format(uuid)
//	@Success		204
//	@Failure		400	{object}	middleware.HttpError
//	@Failure		404	{object}	middleware.HttpError
//	@Failure		500	{object}	middleware.HttpError
//	@Router			/movies/{movieID} [delete]
func MoviesDelete(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	movie := GetContextMovie(c)

	err := models.DeleteMovie(tx, movie.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}
