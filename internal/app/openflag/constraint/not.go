package constraint

import "github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

// NotConstraint represents Openflag not constraint.
type NotConstraint struct{}

// Name is an implementation for the Constraint interface.
func (n NotConstraint) Name() string {
	return NotConstraintName
}

// Validate is an implementation for the Constraint interface.
func (n NotConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (n *NotConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (n NotConstraint) Evaluate(e model.Entity) bool {
	return false
}
