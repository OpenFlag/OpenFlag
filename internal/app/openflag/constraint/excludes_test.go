package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type ExcludesConstraintSuite struct {
	ConstraintSuite
}

func (suite *ExcludesConstraintSuite) TestExcludesConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"]}`,
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
				Name: constraint.ExcludesConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"values": ["10", "11"], "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   8,
				Type: "8",
			},
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
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
			EvaluateExpected: true,
		},
		{
			Name: "successfully create constraint and evaluate 4",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"]}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID: 10,
			},
			EvaluateExpected: false,
		},
		{
			Name: "successfully create constraint and evaluate 5",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"values": ["10", "11"], "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:   8,
				Type: "10",
			},
			EvaluateExpected: false,
		},
		{
			Name: "successfully create constraint and evaluate 6",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
				Parameters: json.RawMessage(
					`{"values": ["10", "11"], "property": "test"}`,
				),
			},
			ErrExpected: false,
			Entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "10"},
			},
			EvaluateExpected: false,
		},
		{
			Name: "failed to create constraint with invalid parameters",
			Constraint: model.Constraint{
				Name: constraint.ExcludesConstraintName,
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

func TestExcludesConstraintSuite(t *testing.T) {
	suite.Run(t, new(ExcludesConstraintSuite))
}
