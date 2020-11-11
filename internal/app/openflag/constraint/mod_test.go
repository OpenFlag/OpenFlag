package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/stretchr/testify/suite"
)

type ModConstraintSuite struct {
	ConstraintSuite
}

func (suite *ModConstraintSuite) TestModConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: constraint.RawConstraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 2}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 10,
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: constraint.RawConstraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`{"value": 2, "property": "%s"}`,
						constraint.EntityTypeProperty,
					),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   9,
				Type: "10",
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: constraint.RawConstraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 2, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      9,
				Type:    "9",
				Context: map[string]string{"test": "8"},
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 4",
			Constraint: constraint.RawConstraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 2, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      9,
				Type:    "9",
				Context: map[string]string{"test": "9"},
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint with invalid parameter",
			Constraint: constraint.RawConstraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 1, "property": "test"}`,
				),
			},
			ErrExpected: true,
			Entity: model.Entity{
				ID: 8,
			},
			EvaluateExpected: false,
		},
	}

	suite.RunCases(cases)
}

func TestModConstraintSuite(t *testing.T) {
	suite.Run(t, new(ModConstraintSuite))
}
