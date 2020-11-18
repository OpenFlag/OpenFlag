package request

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minEntityLen = 1
)

type (
	// Entity represents the context of what we are going to assign the variant on.
	// Usually, OpenFlag expects the context coming with the entity,
	// so that one can define constraints based on the context of the entity.
	Entity struct {
		EntityID      int64             `json:"entity_id"`
		EntityType    string            `json:"entity_type"`
		EntityContext map[string]string `json:"entity_context,omitempty"`
	}

	// EvaluationRequest represents a request for evaluation of some entities.
	EvaluationRequest struct {
		Entities          []Entity `json:"entities"`
		Flags             []string `json:"flags,omitempty"`
		SaveContexts      bool     `json:"save_contexts,omitempty"`
		UseStoredContexts bool     `json:"use_stored_contexts,omitempty"`
	}
)

// Validate validates EvaluationRequest struct.
func (e EvaluationRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(
			&e.Entities,
			validation.Required,
			validation.Length(minEntityLen, 0),
		),
		validation.Field(
			&e.Flags,
			validation.By(func(value interface{}) error {
				for _, flag := range e.Flags {
					if !nameRegex.MatchString(flag) {
						return errors.New("invalid flag")
					}
				}

				return nil
			}),
		),
	)
}
