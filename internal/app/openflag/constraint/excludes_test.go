package constraint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
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
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID: 8,
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID: 10,
					},
					ResultExpected: false,
				},
			},
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
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:   8,
						EntityType: "8",
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID:   8,
						EntityType: "10",
					},
					ResultExpected: false,
				},
			},
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
			Evaluations: []struct {
				Entity         model.Entity
				ResultExpected bool
			}{
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "8"},
					},
					ResultExpected: true,
				},
				{
					Entity: model.Entity{
						EntityID:      8,
						EntityType:    "t",
						EntityContext: map[string]string{"test": "10"},
					},
					ResultExpected: false,
				},
			},
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
		},
	}

	suite.RunCases(cases)
}

func TestExcludesConstraintSuite(t *testing.T) {
	suite.Run(t, new(ExcludesConstraintSuite))
}
