package constraint

import (
	"strconv"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/sirupsen/logrus"
)

// LessThanEqualConstraint represents Openflag less than equal constraint.
type LessThanEqualConstraint struct {
	Value    float64 `json:"value"`
	Property string  `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Name() string {
	return LessThanEqualConstraintName
}

// Validate is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Validate() error {
	return nil
}

// Initialize is an implementation for the Constraint interface.
func (l *LessThanEqualConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (l LessThanEqualConstraint) Evaluate(e model.Entity) bool {
	property, ok := GetProperty(l.Property, e)
	if !ok {
		return false
	}

	value, err := strconv.ParseFloat(property, 64)
	if err != nil {
		logrus.Errorf(
			"invalid property for less than equal constraint => property: %s, value: %s, err: %s",
			l.Property, property, err.Error(),
		)

		return false
	}

	return value <= l.Value
}
