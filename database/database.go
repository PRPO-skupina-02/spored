package database

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/orgs/PRPO-skupina-02/Spored/common"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func OpenAndMigrate() (*gorm.DB, error) {
	db, err := Open()
	if err != nil {
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetProdDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		common.GetEnv("POSTGRES_IP"),
		common.GetEnv("POSTGRES_USERNAME"),
		common.GetEnv("POSTGRES_PASSWORD"),
		common.GetEnv("POSTGRES_DATABASE_NAME"),
		common.GetEnv("POSTGRES_PORT"))
}

func Open() (*gorm.DB, error) {
	dsn := GetProdDSN()

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	slog.Debug("Successfully connected to database", "dsn", dsn)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	instance, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := migratepostgres.WithInstance(instance, &migratepostgres.Config{})
	if err != nil {
		return err
	}

	migrationsFSDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	migrations, err := migrate.NewWithInstance("migrations", migrationsFSDriver, "postgres", driver)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	slog.Debug("Database migrated successfully")
	return nil
}

func GetMigrationsFs() embed.FS {
	return migrationsFS
}
