package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

const migrationPath = "file://migrations"

func runMigration(logger *logrus.Logger, config *AppConfig) error {
	// postgres://user:password@host:port/database?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		logger.WithError(err).Error("error creating migration")
		return err
	}

	defer m.Close()

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("No new migrations")
		return nil
	}
	if err != nil {
		logger.WithError(err).Error("error running migration")
		return err
	}

	logger.Info("Successfully executed migrations")
	return nil
}
