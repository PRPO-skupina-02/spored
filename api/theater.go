package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orgs/PRPO-skupina-02/Spored/models"
)

type TheaterResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func newTheaterResponse(theater models.Theater) TheaterResponse {
	return TheaterResponse{
		UUID:      theater.UUID,
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
//	@Success		200	{object}	[]TheaterRequest
//	@Failure		400	{object}	HttpError
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/theaters [get]
func TheatersList(c *gin.Context) {
	tx := GetContextTransaction(c)

	var theaters []models.Theater
	if err := tx.Find(&theaters).Error; err != nil {
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
//	@Success		200		{object}	[]TheaterRequest
//	@Failure		400		{object}	HttpError
//	@Failure		404		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/theaters [post]
func TheatersCreate(c *gin.Context) {
	tx := GetContextTransaction(c)

	var req TheaterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	theater := models.Theater{
		UUID: uuid.New(),
		Name: req.Name,
	}
	if err := tx.Create(&theater).Error; err != nil {
		_ = c.Error(err)
		return
	}

	response := newTheaterResponse(theater)

	c.JSON(http.StatusOK, response)
}
