package constraint

import (
	"fmt"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ExcludesConstraint represents Openflag excludes constraint.
type ExcludesConstraint struct {
	valueMap map[string]struct{}
	Values   []string `json:"values"`
	Property string   `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Name() string {
	return ExcludesConstraintName
}

// Validate is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Validate() error {
	return validation.ValidateStruct(&ex,
		validation.Field(
			&ex.Values,
			validation.Required,
			validation.Length(minValueLen, 0),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (ex *ExcludesConstraint) Initialize() error {
	valueMap := make(map[string]struct{})

	for _, value := range ex.Values {
		valueMap[value] = struct{}{}
	}

	ex.valueMap = valueMap

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Evaluate(e model.Entity) bool {
	var property string

	switch ex.Property {
	case "":
		property = fmt.Sprintf("%d", e.ID)
	case EntityTypeProperty:
		property = e.Type
	default:
		var ok bool

		property, ok = e.Context[ex.Property]
		if !ok {
			return false
		}
	}

	_, ok := ex.valueMap[property]

	return !ok
}
