package postgres

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Client) runMigrations() {
	if !c.config.AutoMigrate {
		return
	}

	c.logger.Info("running migrations")

	if _, err := c.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, c.config.Schema)); err != nil {
		panic(fmt.Errorf("failed to create schema: %v", err))
	}

	driver, err := postgres.WithInstance(c.DB(), &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to create postgres driver: %v", err))
	}

	_, file, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(file), "..", "..")

	migrationDir := fmt.Sprintf("file://%s/%s", basePath, c.config.MigrationDir)
	m, err := migrate.NewWithDatabaseInstance(migrationDir, c.config.Database, driver)

	if err != nil {
		panic(fmt.Errorf("failed to create migrate instance: %v", err))
	}

	if err := m.Up(); err != migrate.ErrNoChange && err != migrate.ErrNilVersion && err != nil {
		panic(fmt.Errorf("failed to migrate: %v", err))
	}

	c.logger.Info("migrations completed")
}
