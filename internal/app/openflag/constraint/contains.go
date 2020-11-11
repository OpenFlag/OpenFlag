package constraint

import (
	"fmt"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minValueLen = 1
)

// ContainsConstraint represents Openflag contains constraint.
type ContainsConstraint struct {
	valueMap map[string]struct{}
	Values   []string `json:"values"`
	Property string   `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (c ContainsConstraint) Name() string {
	return ContainsConstraintName
}

// Validate is an implementation for the Constraint interface.
func (c ContainsConstraint) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Values,
			validation.Required,
			validation.Length(minValueLen, 0),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (c *ContainsConstraint) Initialize() error {
	valueMap := make(map[string]struct{})

	for _, value := range c.Values {
		valueMap[value] = struct{}{}
	}

	c.valueMap = valueMap

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (c ContainsConstraint) Evaluate(e model.Entity) bool {
	var property string

	switch c.Property {
	case "":
		property = fmt.Sprintf("%d", e.ID)
	case EntityTypeProperty:
		property = e.Type
	default:
		var ok bool

		property, ok = e.Context[c.Property]
		if !ok {
			return false
		}
	}

	_, ok := c.valueMap[property]

	return ok
}
