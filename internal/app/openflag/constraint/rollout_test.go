package constraint_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type RolloutConstraintSuite struct {
	ConstraintSuite
}

func (suite *RolloutConstraintSuite) TestRolloutConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.RolloutConstraintName,
				Parameters: json.RawMessage(
					`{"lower_bound": 10, "upper_bound": 20}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 15,
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.RolloutConstraintName,
				Parameters: json.RawMessage(
					`{"lower_bound": 10, "upper_bound": 20}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 5,
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint with invalid parameters",
			Constraint: model.Constraint{
				Name: constraint.RolloutConstraintName,
				Parameters: json.RawMessage(
					`{"lower_bound": 20, "upper_bound": 10}`,
				),
			},
			ErrExpected: true,
			Entity: model.Entity{
				ID: 5,
			},
			EvaluateExpected: false,
		},
	}

	suite.RunCases(cases)
}

func TestRolloutConstraintSuite(t *testing.T) {
	suite.Run(t, new(RolloutConstraintSuite))
}
