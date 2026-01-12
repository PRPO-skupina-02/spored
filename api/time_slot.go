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

type TimeSlotResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	RoomID    uuid.UUID `json:"room_id"`
	MovieID   uuid.UUID `json:"movie_id"`
}

func newTimeSlotResponse(timeSlot models.TimeSlot) TimeSlotResponse {
	return TimeSlotResponse{
		ID:        timeSlot.ID,
		CreatedAt: timeSlot.CreatedAt,
		UpdatedAt: timeSlot.UpdatedAt,
		StartTime: timeSlot.StartTime,
		EndTime:   timeSlot.EndTime,
		RoomID:    timeSlot.RoomID,
		MovieID:   timeSlot.MovieID,
	}
}

// TimeSlotsList
//
//	@Id				TimeSlotsList
//	@Summary		List time slots
//	@Description	List time slots
//	@Tags			timeslots
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string	true	"Theater ID"					Format(uuid)
//	@Param			roomID		path		string	true	"Room ID"						Format(uuid)
//	@Param			limit		query		int		false	"Limit the number of responses"	Default(10)
//	@Param			offset		query		int		false	"Offset the first response"		Default(0)
//	@Param			sort		query		string	false	"Sort results"
//	@Success		200			{object}	request.PaginatedResponse{data=[]TimeSlotResponse}
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms/{roomID}/timeslots [get]
func TimeSlotsList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	theater := GetContextTheater(c)
	pagination := request.GetNormalizedPaginationArgs(c)
	sort := request.GetSortOptions(c)

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

	timeSlots, total, err := models.GetRoomTimeSlots(tx, room.ID, pagination, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []TimeSlotResponse{}

	for _, timeSlot := range timeSlots {
		response = append(response, newTimeSlotResponse(timeSlot))
	}

	request.RenderPaginatedResponse(c, response, total)
}

// TimeSlotsShow
//
//	@Id				TimeSlotsShow
//	@Summary		Show time slot
//	@Description	Show time slot
//	@Tags			timeslots
//	@Accept			json
//	@Produce		json
//	@Param			theaterID	path		string	true	"Theater ID"	Format(uuid)
//	@Param			roomID		path		string	true	"Room ID"		Format(uuid)
//	@Param			timeSlotID	path		string	true	"TimeSlot ID"	Format(uuid)
//	@Success		200			{object}	TimeSlotResponse
//	@Failure		400			{object}	middleware.HttpError
//	@Failure		404			{object}	middleware.HttpError
//	@Failure		500			{object}	middleware.HttpError
//	@Router			/theaters/{theaterID}/rooms/{roomID}/timeslots/{timeSlotID} [get]
func TimeSlotsShow(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	roomID, err := request.GetUUIDParam(c, "roomID")
	if err != nil {
		_ = c.Error(err)
		return
	}
	timeSlotID, err := request.GetUUIDParam(c, "timeSlotID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	timeSlot, err := models.GetTimeSlot(tx, roomID, timeSlotID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newTimeSlotResponse(timeSlot))
}
