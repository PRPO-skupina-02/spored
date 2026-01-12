package models

import (
	"log/slog"
	"math"
	"slices"
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

func GetTheaterRooms(tx *gorm.DB, theaterID uuid.UUID, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Room, int, error) {
	var rooms []Room

	query := tx.Model(&Room{}).Where("rooms.theater_id = ?", theaterID).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort), PreloadOrderedTimeSlotsScope).Find(&rooms).Error; err != nil {
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

	if err := tx.Where(&room).Scopes(PreloadOrderedTimeSlotsScope).First(&room).Error; err != nil {
		return room, err
	}

	return room, nil
}

func DeleteRoom(tx *gorm.DB, theaterID, id uuid.UUID) error {
	room := Room{
		ID:        id,
		TheaterID: theaterID,
	}

	if err := tx.Where(&room).Scopes(PreloadOrderedTimeSlotsScope).First(&room).Error; err != nil {
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

func PreloadOrderedTimeSlotsScope(db *gorm.DB) *gorm.DB {
	return db.Preload("TimeSlots", func(db *gorm.DB) *gorm.DB {
		return db.Order("time_slots.start_time")
	})
}

const durationDay = time.Hour * 24

func (r *Room) GetTimes(day time.Time) (openingTime time.Time, closingTime time.Time) {
	baseDayTime := day.Truncate(durationDay)
	openingTime = baseDayTime.Add(time.Hour * time.Duration(r.OpeningHour))
	closingTime = baseDayTime.Add(time.Hour * time.Duration(r.ClosingHour))
	return
}

type TimeSlotGap struct {
	room  *Room
	start time.Time
	end   time.Time
}

func (r *Room) GetTimeSlotGapsForDay(day time.Time) []TimeSlotGap {
	startTime, closingTime := r.GetTimes(day)
	gaps := []TimeSlotGap{}

	for _, timeslot := range r.TimeSlots {
		if startTime.After(closingTime) {
			break
		}

		if !timeslot.CoversInstant(startTime) {
			gaps = append(gaps, TimeSlotGap{
				room:  r,
				start: startTime,
				end:   timeslot.StartTime,
			})
		}

		startTime = timeslot.EndTime
	}

	gaps = append(gaps, TimeSlotGap{
		room:  r,
		start: startTime,
		end:   closingTime,
	})

	return gaps
}

func (tsg *TimeSlotGap) Populate(tx *gorm.DB, movies []Movie) error {
	slog.Debug("Populating time gap", "start", tsg.start, "end", tsg.end)
	startTime := tsg.start
	for startTime.Before(tsg.end) {
		remainingMinutes := int(math.Floor(tsg.end.Sub(startTime).Minutes()))

		possibleMovies := slices.Collect(func(yield func(Movie) bool) {
			for _, movie := range movies {
				if !movie.Active {
					continue
				}
				if movie.LengthMinutes <= remainingMinutes {
					if !yield(movie) {
						return
					}
				}
			}
		})

		slog.Debug("Selecting movie", "startTime", startTime, "optionsLen", len(possibleMovies))

		if len(possibleMovies) == 0 {
			slog.Debug("No more possible movies")
			break
		}

		selectedMovie := WeighedSelectMovie(possibleMovies)
		slog.Debug("Selected movie", "title", selectedMovie.Title)
		calculatedEndTime := selectedMovie.CalculateEndTime(startTime)

		timeSlot := TimeSlot{
			ID:        uuid.New(),
			StartTime: startTime,
			EndTime:   calculatedEndTime,
			RoomID:    tsg.room.ID,
			MovieID:   selectedMovie.ID,
		}

		err := timeSlot.Create(tx)
		if err != nil {
			return err
		}

		startTime = calculatedEndTime
	}

	slog.Debug("Finished populating time gap", "start", tsg.start, "end", tsg.end)
	return nil
}

func (r *Room) PopulateRoom(tx *gorm.DB, now time.Time, days int, movies []Movie) error {
	for day := range days {
		slog.Debug("Refreshing timeslots", "room", r.ID, "day", day)
		baseDayTime := now.Add(durationDay * time.Duration(day))
		gaps := r.GetTimeSlotGapsForDay(baseDayTime)
		for _, gap := range gaps {
			err := gap.Populate(tx, movies)
			if err != nil {
				return err
			}
		}
		slog.Debug("Finished refreshing timeslots", "room", r.ID, "day", day)
	}

	return nil
}

func (r *Room) PruneRoom(tx *gorm.DB, before time.Time) error {
	slog.Debug("Pruning timeslots", "room", r.ID)

	timeslots, _, err := GetRoomTimeSlots(tx, r.ID, nil, &request.SortOptions{Column: "end_time", Desc: false}, nil)
	if err != nil {
		return err
	}

	for _, timeslot := range timeslots {
		if timeslot.EndTime.Before(before) {
			err := DeleteTimeSlot(tx, r.ID, timeslot.ID)
			if err != nil {
				return err
			}
		}
	}

	slog.Debug("Finished pruning timeslots", "room", r.ID)

	return nil
}
