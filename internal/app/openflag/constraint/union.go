package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// UnionConstraint represents Openflag union constraint.
type UnionConstraint struct{}

// Name is an implementation for the Constraint interface.
func (u UnionConstraint) Name() string {
	return UnionConstraintName
}

// Validate is an implementation for the Constraint interface.
func (u UnionConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (u *UnionConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (u UnionConstraint) Evaluate(e model.Entity) bool {
	return false
}
