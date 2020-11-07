package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"    // MySQL driver should have blank import
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver should have blank import
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Imported for its side effects
)

// Migrate provides a method for database migration.
func Migrate(driver string, connStr string, path string) error {
	m, err := migrate.New("file://"+path+"/"+driver, connStr)
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
