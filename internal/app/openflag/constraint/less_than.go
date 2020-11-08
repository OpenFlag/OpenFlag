package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// LessThanConstraint represents Openflag less than constraint.
type LessThanConstraint struct{}

// Name is an implementation for the Constraint interface.
func (l LessThanConstraint) Name() string {
	return LessThanConstraintName
}

// Validate is an implementation for the Constraint interface.
func (l LessThanConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (l *LessThanConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (l LessThanConstraint) Evaluate(e model.Entity) bool {
	return false
}
