package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver should have blank import
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// Migrate provides a method for database migration.
func Migrate(source *bindata.AssetSource, connStr string) error {
	driver, err := bindata.WithInstance(source)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", driver, connStr)
	if err != nil {
		return err
	}

	if err := m.Up(); err == migrate.ErrNoChange {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
