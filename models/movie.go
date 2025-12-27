package models

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Movie struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Title         string
	Description   string
	ImageURL      string
	Rating        float64
	LengthMinutes int
	Active        bool
}

func roundToPrecision(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (m *Movie) BeforeSave(tx *gorm.DB) error {
	m.Rating = roundToPrecision(m.Rating, 1)

	m.Rating = math.Max(math.Min(m.Rating, 10), 0)

	return nil
}

func (m *Movie) Create(tx *gorm.DB) error {
	if err := tx.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Movie) Save(tx *gorm.DB) error {
	if err := tx.Save(m).Error; err != nil {
		return err
	}
	return nil
}

func GetMovies(tx *gorm.DB, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Movie, int, error) {
	var movies []Movie

	query := tx.Model(&Movie{}).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort)).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return movies, int(total), nil
}

func GetMovie(tx *gorm.DB, id uuid.UUID) (Movie, error) {
	movie := Movie{
		ID: id,
	}

	if err := tx.Where(&movie).First(&movie).Error; err != nil {
		return movie, err
	}

	return movie, nil
}

func DeleteMovie(tx *gorm.DB, id uuid.UUID) error {
	movie := Movie{
		ID: id,
	}

	if err := tx.Where(&movie).First(&movie).Error; err != nil {
		return err
	}

	if err := tx.Delete(&movie).Error; err != nil {
		return err
	}
	return nil
}

const breakLengthMinutes = 5

func (m *Movie) CalculateEndTime(startTime time.Time) time.Time {
	totalLength := m.LengthMinutes + breakLengthMinutes
	roundedTo10Mins := math.Ceil(float64(totalLength)/10) * 10
	return startTime.Add(time.Duration(roundedTo10Mins) * time.Minute)
}

func WeighedSelectMovie(movies []Movie) Movie {
	// TODO: weigh random selection
	selection := rand.IntN(len(movies))
	return movies[selection]
}
