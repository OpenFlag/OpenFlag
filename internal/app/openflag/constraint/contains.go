package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// ContainsConstraint represents Openflag contains constraint.
type ContainsConstraint struct{}

// Name is an implementation for the Constraint interface.
func (c ContainsConstraint) Name() string {
	return ContainsConstraintName
}

// Validate is an implementation for the Constraint interface.
func (c ContainsConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (c *ContainsConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (c ContainsConstraint) Evaluate(e model.Entity) bool {
	return false
}
