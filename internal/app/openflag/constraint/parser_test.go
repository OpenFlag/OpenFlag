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
		name        string
		cExp        string
		cMap        map[string]model.Constraint
		errExpected bool
		evaluations []struct {
			entity         model.Entity
			resultExpected bool
		}
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 9,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: false,
				},
			},
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 8,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 5,
					},
					resultExpected: false,
				},
			},
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 8,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 7,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 5,
					},
					resultExpected: false,
				},
			},
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 7,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 8,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 5,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 12,
					},
					resultExpected: false,
				},
			},
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 8,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 5,
					},
					resultExpected: true,
				},
			},
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
			evaluations: []struct {
				entity         model.Entity
				resultExpected bool
			}{
				{
					entity: model.Entity{
						ID: 11,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 12,
					},
					resultExpected: false,
				},
				{
					entity: model.Entity{
						ID: 5,
					},
					resultExpected: true,
				},
				{
					entity: model.Entity{
						ID: 110,
					},
					resultExpected: true,
				},
			},
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

			for _, ev := range tc.evaluations {
				result := co.Evaluate(ev.entity)
				suite.Equal(ev.resultExpected, result)
			}
		})
	}
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
