package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	flagName = "sql_flag"
)

type (
	// Variant represents the possible variation of a flag. For example, control/treatment, green/yellow/red, etc.
	// VariantAttachment represents the dynamic configuration of a variant. For example,
	// if you have a variant for the green button,
	// you can dynamically control what's the hex color of green you want to use (e.g. {"hex_color": "#42b983"}).
	Variant struct {
		Key        string          `json:"key"`
		Attachment json.RawMessage `json:"attachment,omitempty"`
	}

	// Constraint represents rules that we can use to define the audience of the segment.
	// In other words, the audience in the segment is defined by a set of constraints.
	Constraint struct {
		Name       string          `json:"name"`
		Parameters json.RawMessage `json:"parameters"`
	}

	// Segment represents the segmentation, i.e. the set of audience we want to target.
	Segment struct {
		Description     string     `json:"description,omitempty"`
		Constraint      Constraint `json:"constraint"`
		Variants        []Variant  `json:"variants"`
		DefaultVariants []Variant  `json:"default_variants,omitempty"`
	}

	// Flag represents each row of flags table in SQL database.
	Flag struct {
		ID          int64      `gorm:"primary_key"`
		Tags        string     `json:"tags"`
		Description string     `json:"description"`
		Flag        string     `json:"flag"`
		Segments    string     `json:"segments"`
		CreatedAt   time.Time  `json:"created_at"`
		DeletedAt   *time.Time `json:"-"`
	}
)

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
