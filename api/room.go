package api

import (
	"net/http"
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orgs/PRPO-skupina-02/spored/models"
)

type RoomResponse struct {
	ID            uuid.UUID                `json:"id"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Name          string                   `json:"name"`
	Rows          int                      `json:"rows"`
	Columns       int                      `json:"columns"`
	OperatingMode models.RoomOperatingMode `json:"operating_mode"`
	OpeningHour   int                      `json:"opening_hour"`
	ClosingHour   int                      `json:"closing_hour"`
}

func newRoomResponse(room models.Room) RoomResponse {
	return RoomResponse{
		ID:            room.ID,
		CreatedAt:     room.CreatedAt,
		UpdatedAt:     room.UpdatedAt,
		Name:          room.Name,
		Rows:          room.Rows,
		Columns:       room.Columns,
		OperatingMode: room.OperatingMode,
		OpeningHour:   room.OpeningHour,
		ClosingHour:   room.ClosingHour,
	}
}

// RoomsList
//
//	@Id				RoomsList
//	@Summary		List rooms
//	@Description	List rooms
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string	true	"Theater ID"					Format(uuid)
//	@Param			limit		query		int		false	"Limit the number of responses"	Default(10)
//	@Param			offset		query		int		false	"Offset the first response"		Default(0)
//	@Param			sort		query		string	false	"Sort results"
//	@Success		200			{object}	[]RoomResponse
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms [get]
func RoomsList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)
	pagination := request.GetNormalizedPaginationArgs(c)
	sort := request.GetSortOptions(c)

	rooms, total, err := models.GetTheaterRooms(tx, theater.ID, pagination, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []RoomResponse{}

	for _, room := range rooms {
		response = append(response, newRoomResponse(room))
	}

	request.RenderPaginatedResponse(c, response, total)
}

type RoomRequest struct {
	Name          string                   `json:"name" binding:"required,min=3"`
	Rows          int                      `json:"rows" binding:"required,min=1,max=100"`
	Columns       int                      `json:"columns" binding:"required,min=1,max=100"`
	OperatingMode models.RoomOperatingMode `json:"operating_mode" binding:"required,oneof=CLOSED WEEKDAYS WEEKENDS ALL"`
	OpeningHour   int                      `json:"opening_hour" binding:"required,min=0,max=24"`
	ClosingHour   int                      `json:"closing_hour" binding:"required,min=0,max=24"`
}

// RoomsCreate
//
//	@Id				RoomsCreate
//	@Summary		Create room
//	@Description	Create room
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string		true	"Theater ID"	Format(uuid)
//	@Param			request		body		RoomRequest	true	"request body"
//	@Success		200			{object}	RoomResponse
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms [post]
func RoomsCreate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)

	var req RoomRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	room := models.Room{
		ID:            uuid.New(),
		TheaterID:     theater.ID,
		Name:          req.Name,
		Rows:          req.Rows,
		Columns:       req.Columns,
		OperatingMode: req.OperatingMode,
		OpeningHour:   req.OpeningHour,
		ClosingHour:   req.ClosingHour,
	}

	err = room.Create(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, newRoomResponse(room))
}

// RoomsShow
//
//	@Id				RoomsShow
//	@Summary		Show room
//	@Description	Show room
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string	true	"Theater ID"	Format(uuid)
//	@Param			roomID		path		string	true	"Room ID"		Format(uuid)
//	@Success		200			{object}	RoomResponse
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms/{roomID} [get]
func RoomsShow(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)
	id, err := request.GetUUIDParam(c, "roomID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	room, err := models.GetRoom(tx, theater.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newRoomResponse(room))
}

// RoomsUpdate
//
//	@Id				RoomsUpdate
//	@Summary		Update room
//	@Description	Update room
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string		true	"Theater ID"	Format(uuid)
//	@Param			roomID		path		string		true	"Room ID"		Format(uuid)
//	@Param			request		body		RoomRequest	true	"request body"
//	@Success		200			{object}	RoomResponse
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms/{roomID} [put]
func RoomsUpdate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)
	id, err := request.GetUUIDParam(c, "roomID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req RoomRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	room, err := models.GetRoom(tx, theater.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	room.Name = req.Name
	room.Rows = req.Rows
	room.Columns = req.Columns
	room.OperatingMode = req.OperatingMode
	room.OpeningHour = req.OpeningHour
	room.ClosingHour = req.ClosingHour

	err = room.Save(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newRoomResponse(room))
}

// RoomsDelete
//
//	@Id				RoomsDelete
//	@Summary		Delete room
//	@Description	Delete room
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path	string	true	"Theater ID"	Format(uuid)
//	@Param			roomID		path	string	true	"Room ID"		Format(uuid)
//	@Success		204
//	@Failure		400	{object}	middleware.HttpError
//	@Failure		404	{object}	middleware.HttpError
//	@Failure		500	{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms/{roomID} [delete]
func RoomsDelete(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)
	id, err := request.GetUUIDParam(c, "roomID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = models.DeleteRoom(tx, theater.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}
