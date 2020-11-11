package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type ContainsConstraintSuite struct {
	ConstraintSuite
}

func (suite *ContainsConstraintSuite) TestContainsConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate using entity id",
			Constraint: model.Constraint{
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
			Name: "successfully create constraint and evaluate using entity type",
			Constraint: model.Constraint{
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
			Name: "successfully create constraint and evaluate using entity context",
			Constraint: model.Constraint{
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
			Name: "successfully create constraint and evaluate using entity id 2",
			Constraint: model.Constraint{
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
			Name: "successfully create constraint and evaluate using entity type 2",
			Constraint: model.Constraint{
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
			Name: "successfully create constraint and evaluate using entity context 2",
			Constraint: model.Constraint{
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
			Constraint: model.Constraint{
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
