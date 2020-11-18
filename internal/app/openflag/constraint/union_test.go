package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/stretchr/testify/suite"
)

type UnionConstraintSuite struct {
	ConstraintSuite
}

func (suite *UnionConstraintSuite) TestUnionConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.UnionConstraintName,
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
						constraint.LessThanConstraintName,
					),
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID: 8,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID: 5,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID: 11,
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "failed to create constraint (creation of inside constraint)",
			Constraint: model.Constraint{
				Name: constraint.UnionConstraintName,
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
		},
		{
			Name: "failed to create constraint (invalid parameters)",
			Constraint: model.Constraint{
				Name: constraint.UnionConstraintName,
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
								}
							]
						}
					`,
						constraint.LessThanConstraintName,
					),
				),
			},
			ErrExpected: true,
		},
	}

	suite.RunCases(cases)
}

func TestUnionConstraintSuite(t *testing.T) {
	suite.Run(t, new(UnionConstraintSuite))
}
