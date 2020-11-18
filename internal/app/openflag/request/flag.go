package request

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minSegmentLen = 1
	maxLimit      = 100

	nameFormat = `^[a-z0-9]+(?:\.[a-z0-9]+)*$`
)

// nolint:gochecknoglobals
var (
	nameRegex = regexp.MustCompile(nameFormat)
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

	// Flag represents a feature flag, an experiment, or a configuration.
	Flag struct {
		Tags        []string  `json:"tags,omitempty"`
		Description string    `json:"description"`
		Flag        string    `json:"flag"`
		Segments    []Segment `json:"segments"`
	}

	// CreateFlagRequest represents a request body for creating a flag.
	CreateFlagRequest struct {
		Flag
	}

	// UpdateFlagRequest represents a request body for updating a flag.
	UpdateFlagRequest struct {
		Flag
	}

	// FindFlagsByTagRequest represents a request body for finding flags that hav given tag.
	FindFlagsByTagRequest struct {
		Tag string `json:"tag"`
	}

	// FindFlagHistoryRequest represents a request body for finding history of a flag.
	FindFlagHistoryRequest struct {
		Flag string `json:"flag"`
	}

	// FindFlagsRequest represents a request body for finding flags using offset and limit.
	FindFlagsRequest struct {
		Offset    int        `json:"offset"`
		Limit     int        `json:"limit"`
		Timestamp *time.Time `json:"timestamp"`
	}
)

// Validate validates FindFlagsRequest struct.
func (f FindFlagsRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(
			&f.Offset,
			validation.Min(0),
		),
		validation.Field(
			&f.Limit,
			validation.Required,
			validation.Min(0), validation.Max(maxLimit),
		),
	)
}

// Validate validates FindFlagsByTagRequest struct.
func (f FindFlagsByTagRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(
			&f.Tag,
			validation.Required,
			validation.Match(nameRegex),
		),
	)
}

// Validate validates FindFlagHistoryRequest struct.
func (f FindFlagHistoryRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(
			&f.Flag,
			validation.Required,
			validation.Match(nameRegex),
		),
	)
}

// Validate validates Variant struct.
func (v Variant) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(
			&v.VariantKey,
			validation.Required,
			validation.Match(nameRegex),
		),
	)
}

// Validate validates Segment struct.
// nolint:funlen
func (s Segment) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(
			&s.Description,
			validation.Required,
		),
		validation.Field(
			&s.Constraints,
			validation.Required,
			validation.By(func(value interface{}) error {
				for _, c := range s.Constraints {
					find := false

					for _, name := range constraint.BasicConstraints() {
						if c.Name == name {
							find = true
							break
						}
					}

					if !find {
						return errors.New("invalid segment constraints")
					}

					if err := constraint.Validate(c.Name, c.Parameters); err != nil {
						return err
					}
				}

				return nil
			}),
		),
		validation.Field(
			&s.Expression,
			validation.Required,
			validation.By(func(value interface{}) error {
				parser := constraint.Parser{}

				constraints := map[string]model.Constraint{}

				for k, v := range s.Constraints {
					constraints[k] = model.Constraint{
						Name:       v.Name,
						Parameters: v.Parameters,
					}
				}

				c, err := parser.Parse(s.Expression, constraints)
				if err != nil {
					return err
				}

				if err := constraint.Validate(c.Name, c.Parameters); err != nil {
					return err
				}

				return nil
			}),
		),
		validation.Field(
			&s.Variant,
			validation.Required,
		),
	)
}

// Validate validates Flag struct.
func (f Flag) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(
			&f.Description,
			validation.Required,
		),
		validation.Field(
			&f.Flag,
			validation.Required,
			validation.Match(nameRegex),
		),
		validation.Field(
			&f.Segments,
			validation.Required,
			validation.Length(minSegmentLen, 0),
		),
		validation.Field(
			&f.Tags,
			validation.By(func(value interface{}) error {
				for _, tag := range f.Tags {
					if !nameRegex.MatchString(tag) {
						return errors.New("invalid flag tag")
					}
				}

				return nil
			}),
		),
	)
}
