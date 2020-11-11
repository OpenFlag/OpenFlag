package constraint_test

import (
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/suite"
)

type (
	ConstraintSuite struct {
		suite.Suite
	}

	ConstraintTestCase struct {
		Name             string           `json:"name"`
		Constraint       model.Constraint `json:"constraint"`
		ErrExpected      bool             `json:"err_expected"`
		Entity           model.Entity     `json:"entity"`
		EvaluateExpected bool             `json:"evaluate_expected"`
	}
)

func (suite *ConstraintSuite) RunCases(cases []ConstraintTestCase) {
	for i := range cases {
		tc := cases[i]

		suite.Run(tc.Name, func() {
			c, err := constraint.New(tc.Constraint)
			if tc.ErrExpected {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			result := c.Evaluate(tc.Entity)
			suite.Equal(tc.EvaluateExpected, result)
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
