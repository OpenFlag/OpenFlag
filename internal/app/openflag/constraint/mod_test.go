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
			Constraint: model.Constraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 2}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						ID: 10,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID: 11,
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`{"value": 2, "property": "%s"}`,
						constraint.EntityTypeProperty,
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
						ID:   9,
						Type: "10",
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:   9,
						Type: "11",
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 2, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						ID:      9,
						Type:    "9",
						Context: map[string]string{"test": "8"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:      9,
						Type:    "9",
						Context: map[string]string{"test": "11"},
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "failed to create constraint with invalid parameter",
			Constraint: model.Constraint{
				Name: constraint.ModConstraintName,
				Parameters: json.RawMessage(
					`{"value": 1, "property": "test"}`,
				),
			},
			ErrExpected: true,
		},
	}

	suite.RunCases(cases)
}

func TestModConstraintSuite(t *testing.T) {
	suite.Run(t, new(ModConstraintSuite))
}
