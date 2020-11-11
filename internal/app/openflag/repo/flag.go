package repo

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	flagName = "sql_flag"
)

// Flag represents each row of SQL database.
type Flag struct {
	ID          int64 `gorm:"primary_key"`
	Tags        string
	Description string
	Flag        string
	Segments    string
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

// FlagRepo provides an interface for communicating with database.
type FlagRepo interface {
	Create(rule *Flag) error
	Delete(id int64) error
	Update(id int64, rule *Flag) error
	FindAll() ([]Flag, error)
	FindByID(id int64) (*Flag, error)
	FindByTag(tag string) ([]Flag, error)
	FindByFlag(flag string) (*Flag, error)
	FindFlags(offset int, limit int, t time.Time) ([]Flag, error)
}

// SQLFlagRepo is an implementation of FlagRepo for SQL databases.
type SQLFlagRepo struct {
	DB *gorm.DB
}

// Create creates a flag in SQL database.
func (s SQLFlagRepo) Create(rule *Flag) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "create", startTime, finalErr) }()

	return nil
}

// Delete deletes a flag from SQL database.
func (s SQLFlagRepo) Delete(id int64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "delete", startTime, finalErr) }()

	return nil
}

// Update updates a flag in SQL database.
func (s SQLFlagRepo) Update(id int64, rule *Flag) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "update", startTime, finalErr) }()

	return nil
}

// FindAll finds all flags from SQL database.
func (s SQLFlagRepo) FindAll() (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_all", startTime, finalErr) }()

	return nil, nil
}

// FindByID finds a flag with it's given id from SQL database.
func (s SQLFlagRepo) FindByID(id int64) (_ *Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_id", startTime, finalErr) }()

	return nil, nil
}

// FindByTag finds flags with it's given tag from SQL database.
func (s SQLFlagRepo) FindByTag(tag string) (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_tag", startTime, finalErr) }()

	return nil, nil
}

// FindByFlag finds a flag with it's given flag from SQL database.
func (s SQLFlagRepo) FindByFlag(flag string) (_ *Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_flag", startTime, finalErr) }()

	return nil, nil
}

// FindFlags finds flags with given offset and limit from SQL database.
func (s SQLFlagRepo) FindFlags(offset int, limit int, t time.Time) (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_flags", startTime, finalErr) }()

	return nil, nil
}
