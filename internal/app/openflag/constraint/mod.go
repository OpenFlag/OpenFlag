package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// ModConstraint represents Openflag mod constraint.
type ModConstraint struct{}

// Name is an implementation for the Constraint interface.
func (m ModConstraint) Name() string {
	return ModConstraintName
}

// Validate is an implementation for the Constraint interface.
func (m ModConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (m *ModConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (m ModConstraint) Evaluate(e model.Entity) bool {
	return false
}
