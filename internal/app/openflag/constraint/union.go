package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

// UnionConstraint represents Openflag union constraint.
type UnionConstraint struct {
	constraints []Constraint
	Constraints []model.Constraint `json:"constraints"`
}

// Name is an implementation for the Constraint interface.
func (u UnionConstraint) Name() string {
	return UnionConstraintName
}

// Validate is an implementation for the Constraint interface.
func (u UnionConstraint) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.constraints,
			validation.Required,
			validation.Length(minConstraintsLen, 0),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (u *UnionConstraint) Initialize() error {
	for _, raw := range u.Constraints {
		c, err := New(raw.Name, raw.Parameters)
		if err != nil {
			return err
		}

		u.constraints = append(u.constraints, c)
	}

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (u UnionConstraint) Evaluate(e model.Entity) bool {
	for _, c := range u.constraints {
		if c.Evaluate(e) {
			return true
		}
	}

	return false
}
