package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"

	"github.com/stretchr/testify/suite"
)

type BiggerThanEqualConstraintSuite struct {
	ConstraintSuite
}

func (suite *BiggerThanEqualConstraintSuite) TestBiggerThanEqualConstraint() {
	cases := []ConstraintTestCase{
		{
			Name: "successfully create constraint and evaluate 1",
			Constraint: model.Constraint{
				Name: constraint.BiggerThanEqualConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID: 11,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID: 9,
					},
					ResultExpected: false,
				},
				{
					Entity: model.Entity{
						EntityID: 10,
					},
					ResultExpected: true,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 2",
			Constraint: model.Constraint{
				Name: constraint.BiggerThanEqualConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"value": 10, "property": "%s"}`, constraint.EntityTypeProperty),
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:   8,
						EntityType: "11",
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID:   8,
						EntityType: "9",
					},
					ResultExpected: false,
				},
				{
					Entity: model.Entity{
						EntityID:   8,
						EntityType: "10",
					},
					ResultExpected: true,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 3",
			Constraint: model.Constraint{
				Name: constraint.BiggerThanEqualConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "11"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "9"},
					},
					ResultExpected: false,
				},
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "10"},
					},
					ResultExpected: true,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 4",
			Constraint: model.Constraint{
				Name: constraint.BiggerThanEqualConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "1"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "-1"},
					},
					ResultExpected: false,
				},
			},
		},
		{
			Name: "successfully create constraint and evaluate 5",
			Constraint: model.Constraint{
				Name: constraint.BiggerThanEqualConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			ErrExpected: false,
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "t"},
					},
					ResultExpected: false,
				},
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "1"},
					},
					ResultExpected: true,
				},
			},
		},
	}

	suite.RunCases(cases)
}

func TestBiggerThanEqualConstraintSuite(t *testing.T) {
	suite.Run(t, new(BiggerThanEqualConstraintSuite))
}
