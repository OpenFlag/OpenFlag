package constraint_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
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
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID: 15,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID: 5,
					},
					ResultExpected: false,
				},
				{
					Entity: model.Entity{
						EntityID: 25,
					},
					ResultExpected: false,
				},
			},
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
		},
	}

	suite.RunCases(cases)
}

func TestRolloutConstraintSuite(t *testing.T) {
	suite.Run(t, new(RolloutConstraintSuite))
}
