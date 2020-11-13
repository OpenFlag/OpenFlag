package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/stretchr/testify/suite"
)

type MotConstraintSuite struct {
	ConstraintSuite
}

func (suite *MotConstraintSuite) TestMotConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.NotConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraint": {
								"name": "%s",
								"parameters": {
									"value": 10
								}
							}
						}
					`,
						constraint.LessThanConstraintName,
					),
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
				Name: constraint.NotConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraint": {
								"name": "%s",
								"parameters": {
									"value": 10
								}
							}
						}
					`,
						constraint.LessThanConstraintName,
					),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 5,
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint (creation of inside constraint)",
			Constraint: model.Constraint{
				Name: constraint.NotConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraint": {
								"name": "%s",
								"parameters": {
									"value": "10""
								}
							}
						}
					`,
						constraint.LessThanConstraintName,
					),
				),
			},
			ErrExpected: true,
			Entity: model.Entity{
				ID: 15,
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint (invalid parameters)",
			Constraint: model.Constraint{
				Name: constraint.NotConstraintName,
				Parameters: json.RawMessage(
					`{}`,
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

func TestMotConstraintSuite(t *testing.T) {
	suite.Run(t, new(MotConstraintSuite))
}
