package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/stretchr/testify/suite"
)

type MatchConstraintSuite struct {
	ConstraintSuite
}

func (suite *MatchConstraintSuite) TestMatchConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.MatchConstraintName,
				Parameters: json.RawMessage(
					`{"expresion": "^[a-z0-9]+(?:\\.[a-z0-9]+)*$"}`,
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
			},
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.MatchConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(
						`{"expresion": "^[a-z0-9]+(?:\\.[a-z0-9]+)*$", "property": "%s"}`,
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
						ID:   8,
						Type: "hello.how.are.you",
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:   8,
						Type: "hello.how are.you",
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.MatchConstraintName,
				Parameters: json.RawMessage(
					`{"expresion": "^[a-z0-9]+(?:\\.[a-z0-9]+)*$", "property": "test"}`,
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
						Context: map[string]string{"test": "hello.how.are.you"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						ID:      8,
						Type:    "t",
						Context: map[string]string{"test": "hello.how are.you"},
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "failed to create constraint with invalid parameter",
			Constraint: model.Constraint{
				Name: constraint.MatchConstraintName,
				Parameters: json.RawMessage(
					`{"property": "test"}`,
				),
			},
			ErrExpected: true,
		},
	}

	suite.RunCases(cases)
}

func TestMatchConstraintSuite(t *testing.T) {
	suite.Run(t, new(MatchConstraintSuite))
}
