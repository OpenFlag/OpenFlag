package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/stretchr/testify/suite"
)

type ContainsConstraintSuite struct {
	ConstraintSuite
}

func (suite *ContainsConstraintSuite) TestContainsConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"]}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						ID: 11,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID: 9,
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"values": ["10", "11"], "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						ID:   8,
						Type: "11",
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:   8,
						Type: "9",
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"], "property": "test"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						ID:      8,
						Type:    "t",
						Context: map[string]string{"test": "11"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:      8,
						Type:    "t",
						Context: map[string]string{"test": "9"},
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "failed to create constraint with invalid parameters",
			Constraint: model.Constraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": []}`,
				),
			},
			ErrExpected: true,
		},
	}

	suite.RunCases(cases)
}

func TestContainsConstraintSuite(t *testing.T) {
	suite.Run(t, new(ContainsConstraintSuite))
}
