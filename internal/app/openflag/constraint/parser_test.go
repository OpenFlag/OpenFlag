package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

func (suite *ParserSuite) TestParserSuite() {
	cases := []struct {
		name             string
		cExp             string
		cMap             map[string]model.Constraint
		errExpected      bool
		entity           model.Entity
		evaluateExpected bool
	}{
		{
			name: "successfully parse expression 1",
			cExp: "A",
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 10}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 9,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully parse expression 2",
			cExp: fmt.Sprintf("A %s B", constraint.IntersectionConstraintName),
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 10}
					`),
				},
				"B": {
					Name: constraint.BiggerThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 6}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 8,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully parse expression 3",
			cExp: fmt.Sprintf(
				"A %s B %s C",
				constraint.IntersectionConstraintName,
				constraint.IntersectionConstraintName,
			),
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 10}
					`),
				},
				"B": {
					Name: constraint.BiggerThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 6}
					`),
				},
				"C": {
					Name: constraint.ContainsConstraintName,
					Parameters: json.RawMessage(`
						{"values": ["8", "9"]}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 8,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully parse expression 4",
			cExp: fmt.Sprintf(
				"(A %s B) %s C",
				constraint.IntersectionConstraintName,
				constraint.UnionConstraintName,
			),
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 10}
					`),
				},
				"B": {
					Name: constraint.BiggerThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 6}
					`),
				},
				"C": {
					Name: constraint.ContainsConstraintName,
					Parameters: json.RawMessage(`
						{"values": ["8", "9", "11"]}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 11,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully parse expression 5",
			cExp: fmt.Sprintf(
				"%s(A %s B)",
				constraint.NotConstraintName,
				constraint.IntersectionConstraintName,
			),
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 10}
					`),
				},
				"B": {
					Name: constraint.BiggerThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 6}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 11,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully parse expression 6",
			cExp: fmt.Sprintf(
				"(%s(A %s B)) %s C",
				constraint.NotConstraintName,
				constraint.IntersectionConstraintName,
				constraint.UnionConstraintName,
			),
			cMap: map[string]model.Constraint{
				"A": {
					Name: constraint.LessThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 100}
					`),
				},
				"B": {
					Name: constraint.BiggerThanConstraintName,
					Parameters: json.RawMessage(`
						{"value": 6}
					`),
				},
				"C": {
					Name: constraint.ContainsConstraintName,
					Parameters: json.RawMessage(`
						{"values": ["8", "9", "11"]}
					`),
				},
			},
			errExpected: false,
			entity: model.Entity{
				ID: 11,
			},
			evaluateExpected: true,
		},
	}

	parser := constraint.Parser{}

	for i := range cases {
		tc := cases[i]

		suite.Run(tc.name, func() {
			c, err := parser.Parse(tc.cExp, tc.cMap)
			if tc.errExpected {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			co, err := constraint.New(c.Name, c.Parameters)
			suite.NoError(err)

			result := co.Evaluate(tc.entity)
			suite.Equal(tc.evaluateExpected, result)
		})
	}
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
