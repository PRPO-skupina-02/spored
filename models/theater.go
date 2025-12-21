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

func GetTheaters(tx *gorm.DB, offset, limit int, sort *request.SortOptions) ([]Theater, int, error) {
	var theaters []Theater

	query := tx.Scopes(request.PaginateScope(offset, limit))

	if sort != nil {
		query = query.Scopes(request.SortScope(sort))
	}

	if err := query.Find(&theaters).Error; err != nil {
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

	if err := tx.Where(&theater).First(&theater).Error; err != nil {
		return err
	}

	if err := tx.Delete(&theater).Error; err != nil {
		return err
	}
	return nil
}
