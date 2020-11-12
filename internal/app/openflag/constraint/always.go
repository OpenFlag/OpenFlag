package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// AlwaysConstraint represents Openflag always constraint.
type AlwaysConstraint struct{}

// Name is an implementation for the Constraint interface.
func (a AlwaysConstraint) Name() string {
	return AlwaysConstraintName
}

// Validate is an implementation for the Constraint interface.
func (a AlwaysConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (a *AlwaysConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (a AlwaysConstraint) Evaluate(e model.Entity) bool {
	return true
}
