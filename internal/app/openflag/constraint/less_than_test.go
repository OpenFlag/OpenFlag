package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"

	"github.com/stretchr/testify/suite"
)

type LessThanConstraintSuite struct {
	ConstraintSuite
}

func (suite *LessThanConstraintSuite) TestLessThanConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.LessThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 9,
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.LessThanConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"value": 10, "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   11,
				Type: "9",
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.LessThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      11,
				Type:    "t",
				Context: map[string]string{"test": "9"},
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 4",
			Constraint: model.Constraint{
				Name: constraint.LessThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "-1"},
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 5",
			Constraint: model.Constraint{
				Name: constraint.LessThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "t"},
			},
			EvaluateExpected: false,
		},
	}

	suite.RunCases(cases)
}

func TestLessThanConstraintSuite(t *testing.T) {
	suite.Run(t, new(LessThanConstraintSuite))
}
