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
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"]}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 11,
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"values": ["10", "11"], "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   8,
				Type: "11",
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"], "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "11"},
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 4",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"]}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 8,
			},
			EvaluateExpected: false,
		},
		{
			Name: "successfully create constraint and evaluate 5",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"values": ["10", "11"], "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   8,
				Type: "8",
			},
			EvaluateExpected: false,
		},
		{
			Name: "successfully create constraint and evaluate 6",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"], "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "8"},
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint with invalid parameters",
			Constraint: constraint.RawConstraint{
				Name: constraint.ContainsConstraintName,
				Parameters: json.RawMessage(
					`{"values": []}`,
				),
			},
			ErrExpected: true,
			Entity: model.Entity{
				ID: 8,
			},
		},
	}

	suite.RunCases(cases)
}

func TestContainsConstraintSuite(t *testing.T) {
	suite.Run(t, new(ContainsConstraintSuite))
}
