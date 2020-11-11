package constraint

import (
	"strconv"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
)

const (
	minModeValue = 2
)

// ModConstraint represents Openflag mod constraint.
type ModConstraint struct {
	Value    int64  `json:"value"`
	Property string `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (m ModConstraint) Name() string {
	return ModConstraintName
}

// Validate is an implementation for the Constraint interface.
func (m ModConstraint) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(
			&m.Value,
			validation.Required,
			validation.Min(minModeValue),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (m *ModConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (m ModConstraint) Evaluate(e model.Entity) bool {
	property, ok := GetProperty(m.Property, e)
	if !ok {
		return false
	}

	value, err := strconv.ParseInt(property, 0, 64)
	if err != nil {
		logrus.Errorf(
			"invalid property for mod constraint => property: %s, value: %s, err: %s",
			m.Property, property, err.Error(),
		)

		return false
	}

	return value%m.Value == 0
}
