package api

import (
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orgs/PRPO-skupina-02/spored/models"
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
//	@Success		200			{object}	[]TimeSlotResponse
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
