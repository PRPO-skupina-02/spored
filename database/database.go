package database

import (
	"fmt"

	"github.com/orgs/PRPO-skupina-02/Spored/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		common.GetEnv("POSTGRES_IP"),
		common.GetEnv("POSTGRES_USERNAME"),
		common.GetEnv("POSTGRES_PASSWORD"),
		common.GetEnv("POSTGRES_DATABASE_NAME"),
		common.GetEnv("POSTGRES_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return db, nil
}
