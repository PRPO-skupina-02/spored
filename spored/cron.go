package spored

import (
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
)

func SetupCron(db *gorm.DB) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	j, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(PopulateSpored, db),
	)
	if err != nil {
		return err
	}

	s.Start()

	slog.Info("Cron job started", "id", j.ID())
	return nil
}
