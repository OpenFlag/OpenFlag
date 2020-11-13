package constraint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
)

// Represents constraint names.
const (
	AlwaysConstraintName          = "always"
	ContainsConstraintName        = "contains"
	ExcludesConstraintName        = "excludes"
	MatchConstraintName           = "match"
	RandomConstraintName          = "random"
	RolloutConstraintName         = "rollout"
	IntersectionConstraintName    = "∩"
	UnionConstraintName           = "∪"
	LessThanConstraintName        = "<"
	LessThanEqualConstraintName   = "<="
	BiggerThanConstraintName      = ">"
	BiggerThanEqualConstraintName = ">="
	NotConstraintName             = "!"
	ModConstraintName             = "%"
)

const (
	// EntityTypeProperty represents entity type constraint property.
	EntityTypeProperty = "type"
)

var (
	// ErrInvalidConstraintName represents an error for returning when the given name for a constraint not found.
	ErrInvalidConstraintName = errors.New("invalid constraint name")
)

// Constraint represents an interface for defining OpenFlag Constraints.
type Constraint interface {
	// Name returns the constraint name.
	Name() string
	// Validate validates the constraint parameters.
	Validate() error
	// Initialize Initializes the constraint.
	Initialize() error
	// Evaluate evaluates an entity in a constraint.
	Evaluate(model.Entity) bool
}

// BasicConstraints returns basic constraints.
func BasicConstraints() []string {
	return []string{
		AlwaysConstraintName,
		ContainsConstraintName,
		ExcludesConstraintName,
		MatchConstraintName,
		RandomConstraintName,
		RolloutConstraintName,
		LessThanConstraintName,
		LessThanEqualConstraintName,
		BiggerThanConstraintName,
		BiggerThanEqualConstraintName,
		ModConstraintName,
	}
}

// Find finds the constraint using the given name.
func Find(name string) (Constraint, error) {
	switch name {
	case AlwaysConstraintName:
		return &AlwaysConstraint{}, nil
	case ContainsConstraintName:
		return &ContainsConstraint{}, nil
	case ExcludesConstraintName:
		return &ExcludesConstraint{}, nil
	case IntersectionConstraintName:
		return &IntersectionConstraint{}, nil
	case MatchConstraintName:
		return &MatchConstraint{}, nil
	case RandomConstraintName:
		return &RandomConstraint{}, nil
	case RolloutConstraintName:
		return &RolloutConstraint{}, nil
	case UnionConstraintName:
		return &UnionConstraint{}, nil
	case LessThanConstraintName:
		return &LessThanConstraint{}, nil
	case LessThanEqualConstraintName:
		return &LessThanEqualConstraint{}, nil
	case BiggerThanConstraintName:
		return &BiggerThanConstraint{}, nil
	case BiggerThanEqualConstraintName:
		return &BiggerThanEqualConstraint{}, nil
	case NotConstraintName:
		return &NotConstraint{}, nil
	case ModConstraintName:
		return &ModConstraint{}, nil
	default:
		return nil, ErrInvalidConstraintName
	}
}

// Validate validates the constraint using the given name and parameters.
func Validate(name string, parameters json.RawMessage) error {
	c, err := Find(name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(parameters, c); err != nil {
		return err
	}

	if err := c.Initialize(); err != nil {
		return err
	}

	return c.Validate()
}

// New create a new constraint using the given name and parameters.
func New(name string, parameters json.RawMessage) (Constraint, error) {
	c, err := Find(name)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(parameters, c); err != nil {
		return nil, err
	}

	if err := c.Initialize(); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// GetProperty returns the property value of a constraint for apply.
func GetProperty(property string, e model.Entity) (string, bool) {
	var value string

	switch property {
	case "":
		value = fmt.Sprintf("%d", e.ID)
	case EntityTypeProperty:
		value = e.Type
	default:
		var ok bool

		value, ok = e.Context[property]
		if !ok {
			return "", false
		}
	}

	return value, true
}
