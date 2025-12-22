package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string
	Rows    int
	Columns int

	TheaterID uuid.UUID
}

func (r *Room) Create(tx *gorm.DB) error {
	if err := tx.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (r *Room) Save(tx *gorm.DB) error {
	if err := tx.Save(r).Error; err != nil {
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

func DeleteRoom(tx *gorm.DB, id uuid.UUID) error {
	room := Room{
		ID: id,
	}

	if err := tx.Where(&room).First(&room).Error; err != nil {
		return err
	}

	if err := tx.Delete(&room).Error; err != nil {
		return err
	}
	return nil
}
