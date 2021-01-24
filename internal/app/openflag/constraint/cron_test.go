package constraint_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type CronConstraintSuite struct {
	ConstraintSuite
}

func (suite *CronConstraintSuite) TestCronConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate",
			Constraint: model.Constraint{
				Name: constraint.CronConstraintName,
				Parameters: json.RawMessage(
					`{"expression": "* * * * *"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID: 1,
					},
					ResultExpected: true,
				},
			},
		},
		{
			Name: "failed to create constraint with invalid parameters",
			Constraint: model.Constraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"expression": "* * * * * *"}`,
				),
			},
			ErrExpected: true,
		},
	}

	suite.RunCases(cases)
}

func TestCronConstraintSuite(t *testing.T) {
	suite.Run(t, new(CronConstraintSuite))
}
