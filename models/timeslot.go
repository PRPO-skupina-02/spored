package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeSlot struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	StartTime time.Time
	EndTime   time.Time

	RoomID  uuid.UUID
	Room    Room `gorm:"foreignKey:RoomID" json:"-"`
	MovieID uuid.UUID
	Movie   Movie `gorm:"foreignKey:MovieID" json:"-"`
}

func (ts *TimeSlot) Create(tx *gorm.DB) error {
	if err := tx.Create(ts).Error; err != nil {
		return err
	}
	return nil
}

func (ts *TimeSlot) Save(tx *gorm.DB) error {
	if err := tx.Save(ts).Error; err != nil {
		return err
	}
	return nil
}

func GetRoomTimeSlots(tx *gorm.DB, theaterID uuid.UUID, offset, limit int, sort *request.SortOptions) ([]TimeSlot, int, error) {
	var timeSlots []TimeSlot

	query := tx.Model(&TimeSlot{}).Where("rooms.theater_id = ?", theaterID).Session(&gorm.Session{})

	if err := query.Debug().Scopes(request.PaginateScope(offset, limit), request.SortScope(sort)).Preload("Movie").Find(&timeSlots).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return timeSlots, int(total), nil
}

func GetTimeSlot(tx *gorm.DB, roomID, timeSlotID uuid.UUID) (TimeSlot, error) {
	timeSlot := TimeSlot{
		ID:     timeSlotID,
		RoomID: roomID,
	}

	if err := tx.Where(&timeSlot).First(&timeSlot).Error; err != nil {
		return timeSlot, err
	}

	return timeSlot, nil
}

func DeleteTimeSlot(tx *gorm.DB, roomID, timeSlotID uuid.UUID) error {
	timeSlot := TimeSlot{
		ID:     timeSlotID,
		RoomID: roomID,
	}

	if err := tx.Where(&timeSlot).First(&timeSlot).Error; err != nil {
		return err
	}

	if err := tx.Delete(&timeSlot).Error; err != nil {
		return err
	}
	return nil
}
