package constraint

import (
	"errors"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	// ErrInvalidRolloutBound represents an error for returning when lower bound and upper bound be invalid.
	ErrInvalidRolloutBound = errors.New("invalid rollout bound")
)

const (
	minBound      = 0
	maxBound      = 99
	maxPercentage = 100
)

// RolloutConstraint represents Openflag rollout constraint.
type RolloutConstraint struct {
	LowerBound int `json:"lower_bound"`
	UpperBound int `json:"upper_bound"`
}

// Name is an implementation for the Constraint interface.
func (r RolloutConstraint) Name() string {
	return RolloutConstraintName
}

// Validate is an implementation for the Constraint interface.
func (r RolloutConstraint) Validate() error {
	if r.LowerBound >= r.UpperBound {
		return ErrInvalidRolloutBound
	}

	return validation.ValidateStruct(&r,
		validation.Field(
			&r.LowerBound,
			validation.Min(minBound),
			validation.Max(maxBound),
		),
		validation.Field(
			&r.UpperBound,
			validation.Required,
			validation.Min(minBound),
			validation.Max(maxBound),
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (r RolloutConstraint) Initialize() error {
	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (r RolloutConstraint) Evaluate(e model.Entity) bool {
	return (e.EntityID%maxPercentage) >= int64(r.LowerBound) &&
		(e.EntityID%maxPercentage) <= int64(r.UpperBound)
}
