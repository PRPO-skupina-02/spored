package spored

import (
	"log/slog"
	"time"

	"github.com/orgs/PRPO-skupina-02/Spored/models"
	"gorm.io/gorm"
)

func PopulateSpored(db *gorm.DB) {
	tx := db.Begin()

	err := func() error {
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
	}()

	if err == nil {
		slog.Info("Schedule successfully populated")
		tx.Commit()
	} else {
		slog.Error("Failed to populate schedule", "err", err)
		tx.Rollback()
	}
}
