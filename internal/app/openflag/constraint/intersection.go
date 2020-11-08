package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// IntersectionConstraint represents Openflag intersection constraint.
type IntersectionConstraint struct{}

// Name is an implementation for the Constraint interface.
func (i IntersectionConstraint) Name() string {
	return IntersectionConstraintName
}

// Validate is an implementation for the Constraint interface.
func (i IntersectionConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (i *IntersectionConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (i IntersectionConstraint) Evaluate(e model.Entity) bool {
	return false
}
