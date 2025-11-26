package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Theater struct {
	UUID      uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
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

func GetTheaters(tx *gorm.DB) ([]Theater, error) {
	var theaters []Theater

	if err := tx.Find(&theaters).Error; err != nil {
		return nil, err
	}

	return theaters, nil
}

func GetTheater(tx *gorm.DB, uuid uuid.UUID) (Theater, error) {
	theater := Theater{
		UUID: uuid,
	}

	if err := tx.Where(&theater).First(&theater).Error; err != nil {
		return theater, err
	}

	return theater, nil
}

func DeleteTheater(tx *gorm.DB, uuid uuid.UUID) error {
	theater := Theater{
		UUID: uuid,
	}

	if err := tx.Where(&theater).First(&theater).Error; err != nil {
		return err
	}

	if err := tx.Delete(&theater).Error; err != nil {
		return err
	}
	return nil
}
