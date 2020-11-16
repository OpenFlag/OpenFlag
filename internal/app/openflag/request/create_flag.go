package request

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minIDValue    = 1
	minSegmentLen = 1

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
		ID         int             `json:"id"`
		Key        string          `json:"key"`
		Attachment json.RawMessage `json:"attachment"`
	}

	// Constraint represents rules that we can use to define the audience of the segment.
	// In other words, the audience in the segment is defined by a set of constraints.
	Constraint struct {
		Name       string          `json:"name"`
		Parameters json.RawMessage `json:"parameters"`
	}

	// Segment represents the segmentation, i.e. the set of audience we want to target.
	Segment struct {
		ID          int                   `json:"id"`
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
)

// Validate validates Variant struct.
func (v Variant) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(
			&v.ID,
			validation.Required,
			validation.Min(minIDValue),
		),
		validation.Field(
			&v.Key,
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
			&s.ID,
			validation.Required,
			validation.Min(minIDValue),
		),
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
			validation.By(func(value interface{}) error {
				sMap := map[int]struct{}{}
				vMap := map[int]struct{}{}

				for _, s := range f.Segments {
					_, ok := sMap[s.ID]
					if ok {
						return errors.New("duplicate segment id")
					}

					sMap[s.ID] = struct{}{}
				}

				for _, s := range f.Segments {
					_, ok := vMap[s.Variant.ID]
					if ok {
						return errors.New("duplicate variant id")
					}

					vMap[s.Variant.ID] = struct{}{}
				}

				return nil
			}),
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
