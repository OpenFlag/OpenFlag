package constraint_test

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type (
	ConstraintSuite struct {
		suite.Suite
	}

	ConstraintTestCase struct {
		Name             string           `json:"name"`
		Constraint       model.Constraint `json:"constraint"`
		ErrExpected      bool             `json:"err_expected"`
		Entity           model.Entity     `json:"entity"`
		EvaluateExpected bool             `json:"evaluate_expected"`
	}
)

func (suite *ConstraintSuite) RunCases(cases []ConstraintTestCase) {
	for i := range cases {
		tc := cases[i]

		suite.Run(tc.Name, func() {
			c, err := constraint.New(tc.Constraint)
			if tc.ErrExpected {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			result := c.Evaluate(tc.Entity)
			suite.Equal(tc.EvaluateExpected, result)
		})
	}
}
