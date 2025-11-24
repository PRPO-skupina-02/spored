package models

import (
	"time"

	"github.com/google/uuid"
)

type Theater struct {
	UUID      uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
