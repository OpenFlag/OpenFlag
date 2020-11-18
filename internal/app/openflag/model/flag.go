package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	flagName = "sql_flag"
)

var (
	// ErrFlagNotFound represents an error for returning when we can't find a flag with given parameters.
	ErrFlagNotFound = errors.New("flag not found")
	// ErrDuplicateFlagFound represents an error for returning when we find a flag with a given key in flag creation.
	ErrDuplicateFlagFound = errors.New("duplicate flag found")
	// ErrInvalidFlagForUpdate represents an error for returning when we find that a flag is invalid in the update method.
	ErrInvalidFlagForUpdate = errors.New("invalid flag for update")
)

type (
	// Variant represents the possible variation of a flag. For example, control/treatment, green/yellow/red, etc.
	// VariantAttachment represents the dynamic configuration of a variant. For example,
	// if you have a variant for the green button,
	// you can dynamically control what's the hex color of green you want to use (e.g. {"hex_color": "#42b983"}).
	Variant struct {
		VariantKey        string          `json:"variant_key"`
		VariantAttachment json.RawMessage `json:"variant_attachment,omitempty"`
	}

	// Constraint represents rules that we can use to define the audience of the segment.
	// In other words, the audience in the segment is defined by a set of constraints.
	Constraint struct {
		Name       string          `json:"name"`
		Parameters json.RawMessage `json:"parameters,omitempty"`
	}

	// Segment represents the segmentation, i.e. the set of audience we want to target.
	Segment struct {
		Description string                `json:"description"`
		Constraints map[string]Constraint `json:"constraints"`
		Expression  string                `json:"expression"`
		Variant     Variant               `json:"variant"`
	}

	// Flag represents each row of flags table in SQL database.
	Flag struct {
		ID          int64      `json:"id" gorm:"primary_key"`
		Tags        *string    `json:"tags,omitempty"`
		Description string     `json:"description"`
		Flag        string     `json:"flag"`
		Segments    string     `json:"segments"`
		CreatedAt   time.Time  `json:"created_at"`
		DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	}
)

// FlagRepo represents an interface for working with persist flags.
type FlagRepo interface {
	Create(flag *Flag) error
	Delete(id int64) error
	Update(id int64, flag *Flag) error
	FindAll() ([]Flag, error)
	FindByID(id int64) (*Flag, error)
	FindByTag(tag string) ([]Flag, error)
	FindByFlag(flag string) ([]Flag, error)
	FindFlags(offset int, limit int, t time.Time) ([]Flag, error)
}

// SQLFlagRepo is an implementation of FlagRepo for SQL databases.
type SQLFlagRepo struct {
	Driver   string
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

// Create creates a flag in SQL database.
func (s SQLFlagRepo) Create(flag *Flag) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "create", startTime, finalErr) }()

	return s.MasterDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("flag = ?", flag.Flag).Take(&Flag{}).Error
		if err == nil {
			return ErrDuplicateFlagFound
		} else if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		if err := tx.Create(flag).Error; err != nil {
			return err
		}

		return nil
	})
}

// Delete deletes a flag from SQL database.
func (s SQLFlagRepo) Delete(id int64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "delete", startTime, finalErr) }()

	return s.MasterDB.Where("id = ?", id).Delete(&Flag{}).Error
}

// Update updates a flag in SQL database.
func (s SQLFlagRepo) Update(id int64, flag *Flag) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "update", startTime, finalErr) }()

	return s.MasterDB.Transaction(func(tx *gorm.DB) error {
		var f Flag

		if err := tx.Where("id = ?", id).Find(&f).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return ErrFlagNotFound
			}

			return err
		}

		if f.Flag != flag.Flag {
			return ErrInvalidFlagForUpdate
		}

		if err := tx.Where("id = ?", id).Delete(&Flag{}).Error; err != nil {
			return err
		}

		if err := tx.Create(flag).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll finds all flags from SQL database.
func (s SQLFlagRepo) FindAll() (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_all", startTime, finalErr) }()

	var result []Flag

	if err := s.SlaveDB.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// FindByID finds a flag with it's given id from SQL database.
func (s SQLFlagRepo) FindByID(id int64) (_ *Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_id", startTime, finalErr) }()

	var result Flag

	if err := s.SlaveDB.Unscoped().Where("id = ?", id).Find(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrFlagNotFound
		}

		return nil, err
	}

	return &result, nil
}

// FindByTag finds flags with it's given tag from SQL database.
func (s SQLFlagRepo) FindByTag(tag string) (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_tag", startTime, finalErr) }()

	var result []Flag

	var query string

	if s.Driver == "postgres" {
		query = fmt.Sprintf("select * from flags where tags::jsonb ? '%s' and deleted_at is null;", tag)
	}

	if err := s.SlaveDB.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// FindByFlag finds history of a flag with it's given key from SQL database.
func (s SQLFlagRepo) FindByFlag(flag string) (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_by_flag", startTime, finalErr) }()

	var result []Flag

	if err := s.SlaveDB.Unscoped().Where("flag = ?", flag).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// FindFlags finds flags with given offset and limit from SQL database.
func (s SQLFlagRepo) FindFlags(offset int, limit int, t time.Time) (_ []Flag, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(flagName, "find_flags", startTime, finalErr) }()

	var result []Flag

	if err := s.SlaveDB.Where("created_at < ?", t).Order("id desc").Offset(offset).Limit(limit).
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
