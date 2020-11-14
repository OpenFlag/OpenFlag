package constraint_test

import (
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/suite"
)

type (
	ConstraintSuite struct {
		suite.Suite
	}

	ConstraintTestCase struct {
		Name        string
		Constraint  model.Constraint
		ErrExpected bool
		Evaluations []struct {
			Entity         model.Entity
			ResultExpected bool
		}
	}
)

func (suite *ConstraintSuite) RunCases(cases []ConstraintTestCase) {
	for i := range cases {
		tc := cases[i]

		suite.Run(tc.Name, func() {
			rc := tc.Constraint

			c, err := constraint.New(rc.Name, rc.Parameters)
			if tc.ErrExpected {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			for _, ev := range tc.Evaluations {
				result := c.Evaluate(ev.Entity)
				suite.Equal(ev.ResultExpected, result)
			}
		})
	}
}

func TestGetProperty(t *testing.T) {
	e := model.Entity{
		ID:      10,
		Type:    "test",
		Context: map[string]string{"context": "context"},
	}

	property, ok := constraint.GetProperty("", e)
	assert.Equal(t, property, "10")
	assert.Equal(t, ok, true)

	property, ok = constraint.GetProperty(constraint.EntityTypeProperty, e)
	assert.Equal(t, property, "test")
	assert.Equal(t, ok, true)

	property, ok = constraint.GetProperty("context", e)
	assert.Equal(t, property, "context")
	assert.Equal(t, ok, true)

	property, ok = constraint.GetProperty("not.found", e)
	assert.Equal(t, property, "")
	assert.Equal(t, ok, false)
}
