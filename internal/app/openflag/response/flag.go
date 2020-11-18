package response

import (
	"encoding/json"
	"time"
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
		ID          int64      `json:"id"`
		Tags        []string   `json:"tags,omitempty"`
		Description string     `json:"description"`
		Flag        string     `json:"flag"`
		Segments    []Segment  `json:"segments"`
		CreatedAt   time.Time  `json:"created_at"`
		DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	}
)
