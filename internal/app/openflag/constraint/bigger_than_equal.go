package constraint

import (
	"strconv"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/sirupsen/logrus"
)

// BiggerThanEqualConstraint represents Openflag bigger than equal constraint.
type BiggerThanEqualConstraint struct {
	Value    float64 `json:"value"`
	Property string  `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Name() string {
	return BiggerThanEqualConstraintName
}

// Validate is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (b *BiggerThanEqualConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (b BiggerThanEqualConstraint) Evaluate(e model.Entity) bool {
	property, ok := GetProperty(b.Property, e)
	if !ok {
		return false
	}

	value, err := strconv.ParseFloat(property, 64)
	if err != nil {
		logrus.Errorf(
			"invalid property for bigger than equal constraint => property: %s, value: %s, err: %s",
			b.Property, property, err.Error(),
		)

		return false
	}

	return value >= b.Value
}
