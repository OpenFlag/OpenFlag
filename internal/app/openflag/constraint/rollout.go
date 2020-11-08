package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// RolloutConstraint represents Openflag rollout constraint.
type RolloutConstraint struct{}

// Name is an implementation for the Constraint interface.
func (r RolloutConstraint) Name() string {
	return RolloutConstraintName
}

// Validate is an implementation for the Constraint interface.
func (r RolloutConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (r RolloutConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (r RolloutConstraint) Evaluate(e model.Entity) bool {
	return false
}
