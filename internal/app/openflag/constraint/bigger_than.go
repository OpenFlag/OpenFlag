package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// BiggerThanConstraint represents Openflag bigger than constraint.
type BiggerThanConstraint struct{}

// Name is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Name() string {
	return BiggerThanConstraintName
}

// Validate is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (b *BiggerThanConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Evaluate(e model.Entity) bool {
	return false
}
