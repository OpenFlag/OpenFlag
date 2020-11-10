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
	if l.Property == "" {
		return float64(e.ID) <= l.Value
	}

	var property string

	if l.Property == EntityTypeProperty {
		property = e.Type
	} else {
		var ok bool

		property, ok = e.Context[l.Property]
		if !ok {
			return false
		}
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
