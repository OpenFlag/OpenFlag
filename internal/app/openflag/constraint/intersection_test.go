package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type IntersectionConstraintSuite struct {
	ConstraintSuite
}

func (suite *IntersectionConstraintSuite) TestIntersectionConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.IntersectionConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraints": [
								{
									"name": "%s",
									"parameters": {
										"value": 10
									}
								},
								{
									"name": "%s",
									"parameters": {
										"value": 6
									}
								}
							]
						}
					`,
						constraint.LessThanConstraintName,
						constraint.BiggerThanConstraintName,
					),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 8,
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.IntersectionConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraints": [
								{
									"name": "%s",
									"parameters": {
										"value": 10
									}
								},
								{
									"name": "%s",
									"parameters": {
										"value": 6
									}
								}
							]
						}
					`,
						constraint.LessThanConstraintName,
						constraint.BiggerThanConstraintName,
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
				Name: constraint.IntersectionConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraints": [
								{
									"name": "%s",
									"parameters": {
										"value": "10"
									}
								},
								{
									"name": "%s",
									"parameters": {
										"value": 6
									}
								}
							]
						}
					`,
						constraint.LessThanConstraintName,
						constraint.BiggerThanConstraintName,
					),
				),
			},
			ErrExpected: true,
			Entity: model.Entity{
				ID: 5,
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint (invalid parameters)",
			Constraint: model.Constraint{
				Name: constraint.IntersectionConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`
						{
							"constraints": [
								{
									"name": "%s",
									"parameters": {
										"value": "10"
									}
								}
							]
						}
					`,
						constraint.LessThanConstraintName,
					),
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

func TestIntersectionConstraintSuite(t *testing.T) {
	suite.Run(t, new(IntersectionConstraintSuite))
}
