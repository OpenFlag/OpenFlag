package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// LessThanEqualConstraint represents Openflag less than equal constraint.
type LessThanEqualConstraint struct{}

// Name is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Name() string {
	return LessThanEqualConstraintName
}

// Validate is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (l *LessThanEqualConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Evaluate(e model.Entity) bool {
	return false
}
