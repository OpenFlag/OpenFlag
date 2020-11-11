package model

type (
	// Entity represents the context of what we are going to assign the variant on.
	// Usually, OpenFlag expects the context coming with the entity,
	// so that one can define constraints based on the context of the entity.
	Entity struct {
		ID      int64             `json:"id"`
		Type    string            `json:"type"`
		Context map[string]string `json:"context,omitempty"`
	}
)
