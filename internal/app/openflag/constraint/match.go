package constraint

import "github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

// MatchConstraint represents Openflag match constraint.
type MatchConstraint struct{}

// Name is an implementation for the Constraint interface.
func (m MatchConstraint) Name() string {
	return MatchConstraintName
}

// Validate is an implementation for the Constraint interface.
func (m MatchConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (m *MatchConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (m MatchConstraint) Evaluate(e model.Entity) bool {
	return false
}
