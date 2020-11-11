package constraint

import (
	"strconv"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/sirupsen/logrus"
)

// BiggerThanConstraint represents Openflag bigger than constraint.
type BiggerThanConstraint struct {
	Value    float64 `json:"value"`
	Property string  `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Name() string {
	return BiggerThanConstraintName
}

// Validate is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (b *BiggerThanConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (b BiggerThanConstraint) Evaluate(e model.Entity) bool {
	property, ok := GetProperty(b.Property, e)
	if !ok {
		return false
	}

	value, err := strconv.ParseFloat(property, 64)
	if err != nil {
		logrus.Errorf(
			"invalid property for bigger than constraint => property: %s, value: %s, err: %s",
			b.Property, property, err.Error(),
		)

		return false
	}

	return value > b.Value
}
