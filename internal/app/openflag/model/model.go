package model

import (
	"encoding/json"
	"regexp"
)

const (
	// KeyFormat is an acceptable regex format for key, name, and type variables.
	KeyFormat = `^[a-z0-9]+(?:\.[a-z0-9]+)*$`
)

//nolint:gochecknoglobals
var (
	// KeyRegex represents a regex for KeyFormat.
	KeyRegex = regexp.MustCompile(KeyFormat)
)

type (
	// Entity represents the context of what we are going to assign the variant on.
	// Usually, OpenFlag expects the context coming with the entity,
	// so that one can define constraints based on the context of the entity.
	Entity struct {
		ID      int64             `json:"id" validate:"required"`
		Type    string            `json:"type" validate:"required"`
		Context map[string]string `json:"context,omitempty"`
	}

	// Variant represents the possible variation of a flag. For example, control/treatment, green/yellow/red, etc.
	Variant struct {
		Key string `json:"key" validate:"required"`
		// Attachment represents the dynamic configuration of a variant. For example,
		// if you have a variant for the green button,
		// you can dynamically control what's the hex color of green you want to use (e.g. {"hex_color": "#42b983"}).
		Attachment json.RawMessage `json:"attachment,omitempty" validate:"required"`
	}

	// Constraint represents rules that we can use to define the audience of the segment.
	// In other words, the audience in the segment is defined by a set of constraints.
	Constraint struct {
		Name       string          `json:"name" validate:"required"`
		Parameters json.RawMessage `json:"parameters,omitempty"`
	}

	// Segment represents the segmentation, i.e. the set of audience we want to target.
	Segment struct {
		Description string     `json:"description" validate:"required"`
		Constraint  Constraint `json:"constraint" validate:"required"`
		Variants    []Variant  `json:"variants" validate:"required"`
	}

	// Flag represents a feature flag, an experiment, or a configuration.
	Flag struct {
		// Tag is a descriptive label attached to a flag for easy lookup and evaluation.
		Tag            *string   `json:"tag,omitempty"`
		Description    string    `json:"description" validate:"required"`
		Flag           string    `json:"flag" validate:"required"`
		Segments       []Segment `json:"segments" validate:"required"`
		DefaultVariant *Variant  `json:"default_variant,omitempty"`
	}
)
