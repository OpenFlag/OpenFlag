package constraint

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/Pallinder/go-randomdata"
)

// RandomConstraint represents Openflag random constraint.
type RandomConstraint struct{}

// Name is an implementation for the Constraint interface.
func (r RandomConstraint) Name() string {
	return RandomConstraintName
}

// Validate is an implementation for the Constraint interface.
func (r RandomConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (r *RandomConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (r RandomConstraint) Evaluate(e model.Entity) bool {
	return randomdata.Boolean()
}
