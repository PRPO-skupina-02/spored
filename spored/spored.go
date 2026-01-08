package spored

import (
	"log/slog"
	"time"

	"github.com/PRPO-skupina-02/spored/models"
	"gorm.io/gorm"
)

func TimeSlotRefresh(db *gorm.DB) {
	tx := db.Begin()

	err := func() error {
		err := PopulateSpored(tx)
		if err != nil {
			return err
		}

		// err = PruneSpored(tx)
		// if err != nil {
		// 	   return err
		// }

		return nil
	}()

	if err == nil {
		slog.Info("TimeSlots successfully refreshed")
		tx.Commit()
	} else {
		slog.Error("Failed to refresh TimeSlots", "err", err)
		tx.Rollback()
	}
}

func PopulateSpored(tx *gorm.DB) error {
	movies, _, err := models.GetMovies(tx, nil, nil)
	if err != nil {
		return err
	}

	theaters, _, err := models.GetTheaters(tx, nil, nil)
	if err != nil {
		return err
	}

	for _, theater := range theaters {
		err = theater.PopulateTheater(tx, time.Now(), 7, movies)
		if err != nil {
			return err
		}
	}

	return nil
}

func PruneSpored(tx *gorm.DB) error {
	theaters, _, err := models.GetTheaters(tx, nil, nil)
	if err != nil {
		return err
	}

	for _, theater := range theaters {
		err = theater.PruneTheater(tx, time.Now().Truncate(time.Hour*24).Add(time.Hour*24*(-7)))
		if err != nil {
			return err
		}
	}

	return nil
}
