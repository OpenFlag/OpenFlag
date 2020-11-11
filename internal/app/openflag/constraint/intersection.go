package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minConstraintsLen = 2
)

// IntersectionConstraint represents Openflag intersection constraint.
type IntersectionConstraint struct {
	constraints    []Constraint
	RawConstraints []RawConstraint `json:"constraints"`
}

// Name is an implementation for the Constraint interface.
func (i IntersectionConstraint) Name() string {
	return IntersectionConstraintName
}

// Validate is an implementation for the Constraint interface.
func (i IntersectionConstraint) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.constraints,
			validation.Required,
			validation.Length(minConstraintsLen, 0),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (i *IntersectionConstraint) Initialize() error {
	for _, raw := range i.RawConstraints {
		c, err := New(raw.Name, raw.Parameters)
		if err != nil {
			return err
		}

		i.constraints = append(i.constraints, c)
	}

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (i IntersectionConstraint) Evaluate(e model.Entity) bool {
	for _, c := range i.constraints {
		if !c.Evaluate(e) {
			return false
		}
	}

	return true
}
