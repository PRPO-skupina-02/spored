package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Theater struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string
}

func (r *Theater) Create(tx *gorm.DB) error {
	if err := tx.Create(r).Error; err != nil {
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

	if err := tx.Where(&theater).First(&theater).Error; err != nil {
		return err
	}

	if err := tx.Delete(&theater).Error; err != nil {
		return err
	}
	return nil
}
