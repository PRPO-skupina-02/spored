package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomOperatingMode string

const (
	Closed   RoomOperatingMode = "CLOSED"
	Weekdays RoomOperatingMode = "WEEKDAYS"
	Weekends RoomOperatingMode = "WEEKENDS"
	All      RoomOperatingMode = "ALL"
)

type Room struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string
	Rows    int
	Columns int

	OperatingMode RoomOperatingMode
	OpeningHour   int
	ClosingHour   int

	TheaterID uuid.UUID
	Theater   Theater    `gorm:"foreignKey:TheaterID" json:"-"`
	TimeSlots []TimeSlot `gorm:"foreignKey:RoomID" json:"-"`
}

func (ts *Room) Create(tx *gorm.DB) error {
	if err := tx.Create(ts).Error; err != nil {
		return err
	}
	return nil
}

func (ts *Room) Save(tx *gorm.DB) error {
	if err := tx.Save(ts).Error; err != nil {
		return err
	}
	return nil
}

func GetTheaterRooms(tx *gorm.DB, theaterID uuid.UUID, offset, limit int, sort *request.SortOptions) ([]Room, int, error) {
	var rooms []Room

	query := tx.Model(&Room{}).Where("rooms.theater_id = ?", theaterID).Session(&gorm.Session{})

	if err := query.Debug().Scopes(request.PaginateScope(offset, limit), request.SortScope(sort)).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return rooms, int(total), nil
}

func GetRoom(tx *gorm.DB, theaterID, roomID uuid.UUID) (Room, error) {
	room := Room{
		ID:        roomID,
		TheaterID: theaterID,
	}

	if err := tx.Where(&room).First(&room).Error; err != nil {
		return room, err
	}

	return room, nil
}

func DeleteRoom(tx *gorm.DB, theaterID, id uuid.UUID) error {
	room := Room{
		ID:        id,
		TheaterID: theaterID,
	}

	if err := tx.Where(&room).Preload("TimeSlots").First(&room).Error; err != nil {
		return err
	}

	for _, timeslot := range room.TimeSlots {
		err := DeleteTimeSlot(tx, id, timeslot.ID)
		if err != nil {
			return err
		}
	}

	if err := tx.Delete(&room).Error; err != nil {
		return err
	}
	return nil
}
