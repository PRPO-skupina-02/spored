package api

import (
	"net/http"
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orgs/PRPO-skupina-02/Spored/models"
)

type TheaterResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func newTheaterResponse(theater models.Theater) TheaterResponse {
	return TheaterResponse{
		ID:        theater.ID,
		CreatedAt: theater.CreatedAt,
		UpdatedAt: theater.UpdatedAt,
		Name:      theater.Name,
	}
}

// TheatersList
//
//	@Id				TheatersList
//	@Summary		List theaters
//	@Description	List theaters
//	@Tags			theaters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]TheaterResponse
//	@Failure		400	{object}	HttpError
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/theaters [get]
func TheatersList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)

	theaters, err := models.GetTheaters(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []TheaterResponse{}

	for _, theater := range theaters {
		response = append(response, newTheaterResponse(theater))
	}

	c.JSON(http.StatusOK, response)
}

type TheaterRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

// TheatersCreate
//
//	@Id				TheatersCreate
//	@Summary		Create theater
//	@Description	Create theater
//	@Tags			theaters
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TheaterRequest	true	"request body"
//	@Success		200		{object}	TheaterResponse
//	@Failure		400		{object}	HttpError
//	@Failure		404		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/theaters [post]
func TheatersCreate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)

	var req TheaterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	theater := models.Theater{
		ID:   uuid.New(),
		Name: req.Name,
	}

	err = theater.Create(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, newTheaterResponse(theater))
}

// TheatersUpdate
//
//	@Id				TheatersUpdate
//	@Summary		Update theater
//	@Description	Update theater
//	@Tags			theaters
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string			true	"Theater UUID"	Format(uuid)
//	@Param			request	body		TheaterRequest	true	"request body"
//	@Success		200		{object}	TheaterResponse
//	@Failure		400		{object}	HttpError
//	@Failure		404		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/theaters/{uuid} [put]
func TheatersUpdate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	uuidParam, err := request.GetUUIDParam(c, "uuid")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req TheaterRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	theater, err := models.GetTheater(tx, uuidParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	theater.Name = req.Name

	err = theater.Save(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newTheaterResponse(theater))
}

// TheatersDelete
//
//	@Id				TheatersDelete
//	@Summary		Delete theater
//	@Description	Delete theater
//	@Tags			theaters
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Theater UUID"	Format(uuid)
//	@Success		204
//	@Failure		400	{object}	HttpError
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/theaters/{uuid} [delete]
func TheatersDelete(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	uuidParam, err := request.GetUUIDParam(c, "uuid")
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = models.DeleteTheater(tx, uuidParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}
