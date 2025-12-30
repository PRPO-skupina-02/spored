package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Theater struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	Rooms []Room `gorm:"foreignKey:TheaterID" json:"-"`
}

func (t *Theater) Create(tx *gorm.DB) error {
	if err := tx.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *Theater) Save(tx *gorm.DB) error {
	if err := tx.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func GetTheaters(tx *gorm.DB, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Theater, int, error) {
	var theaters []Theater

	query := tx.Model(&Theater{}).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort)).Find(&theaters).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return theaters, int(total), nil
}

func GetTheater(tx *gorm.DB, id uuid.UUID) (Theater, error) {
	theater := Theater{
		ID: id,
	}

	if err := tx.Where(&theater).First(&theater).Error; err != nil {
		return theater, err
	}

	return theater, nil
}

func DeleteTheater(tx *gorm.DB, id uuid.UUID) error {
	theater := Theater{
		ID: id,
	}

	if err := tx.Where(&theater).Preload("Rooms").First(&theater).Error; err != nil {
		return err
	}

	for _, room := range theater.Rooms {
		err := DeleteRoom(tx, id, room.ID)
		if err != nil {
			return err
		}
	}

	if err := tx.Delete(&theater).Error; err != nil {
		return err
	}
	return nil
}

func (t *Theater) PopulateTheater(tx *gorm.DB, now time.Time, days int, movies []Movie) error {
	rooms, _, err := GetTheaterRooms(tx, t.ID, nil, nil)
	if err != nil {
		return err
	}

	for _, room := range rooms {
		err := room.PopulateRoom(tx, now, days, movies)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Theater) PruneTheater(tx *gorm.DB, before time.Time) error {
	rooms, _, err := GetTheaterRooms(tx, t.ID, nil, nil)
	if err != nil {
		return err
	}

	for _, room := range rooms {
		err := room.PruneRoom(tx, before)
		if err != nil {
			return err
		}
	}

	return nil
}
