package model

type (
	// Entity represents the context of what we are going to assign the variant on.
	// Usually, OpenFlag expects the context coming with the entity,
	// so that one can define constraints based on the context of the entity.
	Entity struct {
		EntityID      int64             `json:"entity_id"`
		EntityType    string            `json:"entity_type"`
		EntityContext map[string]string `json:"entity_context,omitempty"`
	}
)
