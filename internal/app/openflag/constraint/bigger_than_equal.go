package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// BiggerThanEqualConstraint represents Openflag bigger than equal constraint.
type BiggerThanEqualConstraint struct{}

// Name is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Name() string {
	return BiggerThanEqualConstraintName
}

// Validate is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (b *BiggerThanEqualConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Evaluate(e model.Entity) bool {
	return false
}
