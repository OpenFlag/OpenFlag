package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// ExcludesConstraint represents Openflag excludes constraint.
type ExcludesConstraint struct{}

// Name is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Name() string {
	return ExcludesConstraintName
}

// Validate is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (ex *ExcludesConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (ex ExcludesConstraint) Evaluate(e model.Entity) bool {
	return false
}
