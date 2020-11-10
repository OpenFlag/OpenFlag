package constraint

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type BiggerThanConstraintSuite struct {
	suite.Suite
}

func (suite *BiggerThanConstraintSuite) TestBiggerThanConstraint() {
	cases := []struct {
		name             string
		constraint       model.Constraint
		errExpected      bool
		entity           model.Entity
		evaluateExpected bool
	}{
		{
			name: "successfully create constraint and evaluate using entity id",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10}`,
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID: 11,
			},
			evaluateExpected: true,
		},
		{
			name: "successfully create constraint and evaluate using entity type",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					fmt.Sprintf(`{"value": 10, "property": "%s"}`, EntityTypeProperty),
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID:   8,
				Type: "11",
			},
			evaluateExpected: true,
		},
		{
			name: "successfully create constraint and evaluate using entity context",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 10, "property": "test"}`,
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "11"},
			},
			evaluateExpected: true,
		},
		{
			name: "successfully create constraint and evaluate using zero value",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "1"},
			},
			evaluateExpected: true,
		},
		{
			name: "successfully create constraint and evaluate using invalid property",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "t"},
			},
			evaluateExpected: false,
		},
		{
			name: "successfully create constraint and evaluate using negative property",
			constraint: model.Constraint{
				Name: BiggerThanConstraintName,
				Parameters: json.RawMessage(
					`{"value": 0, "property": "test"}`,
				),
			},
			errExpected: false,
			entity: model.Entity{
				ID:      8,
				Type:    "t",
				Context: map[string]string{"test": "-1"},
			},
			evaluateExpected: false,
		},
	}

	for i := range cases {
		tc := cases[i]

		suite.Run(tc.name, func() {
			constraint, err := New(tc.constraint)
			if tc.errExpected {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			result := constraint.Evaluate(tc.entity)
			suite.Equal(tc.evaluateExpected, result)
		})
	}
}

func TestBiggerThanConstraintSuite(t *testing.T) {
	suite.Run(t, new(BiggerThanConstraintSuite))
}
