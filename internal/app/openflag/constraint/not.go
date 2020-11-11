package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

// NotConstraint represents Openflag not constraint.
type NotConstraint struct {
	constraint    Constraint
	RawConstraint RawConstraint `json:"constraint"`
}

// Name is an implementation for the Constraint interface.
func (n NotConstraint) Name() string {
	return NotConstraintName
}

// Validate is an implementation for the Constraint interface.
func (n NotConstraint) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(
			&n.constraint,
			validation.Required,
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (n *NotConstraint) Initialize() error {
	c, err := New(n.RawConstraint.Name, n.RawConstraint.Parameters)
	if err != nil {
		return err
	}

	n.constraint = c

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (n NotConstraint) Evaluate(e model.Entity) bool {
	return !n.constraint.Evaluate(e)
}
