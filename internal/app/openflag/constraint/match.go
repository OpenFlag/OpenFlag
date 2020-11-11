package constraint

import (
	"fmt"
	"regexp"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

// MatchConstraint represents Openflag match constraint.
type MatchConstraint struct {
	regex     *regexp.Regexp
	Expresion string `json:"expresion"`
	Property  string `json:"property,omitempty"`
}

// Name is an implementation for the Constraint interface.
func (m MatchConstraint) Name() string {
	return MatchConstraintName
}

// Validate is an implementation for the Constraint interface.
func (m MatchConstraint) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(
			&m.Expresion,
			validation.Required,
		),
	)
}

// Initialize is an implementation for the Constraint interface.
func (m *MatchConstraint) Initialize() error {
	regex, err := regexp.Compile(m.Expresion)
	if err != nil {
		return err
	}

	m.regex = regex

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (m MatchConstraint) Evaluate(e model.Entity) bool {
	var property string

	switch m.Property {
	case "":
		property = fmt.Sprintf("%d", e.ID)
	case EntityTypeProperty:
		property = e.Type
	default:
		var ok bool

		property, ok = e.Context[m.Property]
		if !ok {
			return false
		}
	}

	return m.regex.MatchString(property)
}
