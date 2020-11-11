package constraint_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/assert"
)

func TestRandomRule(t *testing.T) {
	rawConstraint := model.Constraint{
		Name:       constraint.RandomConstraintName,
		Parameters: json.RawMessage(`{}`),
	}

	err := constraint.Validate(rawConstraint)
	assert.NoError(t, err)

	c, err := constraint.New(rawConstraint)
	assert.NoError(t, err)

	tr := 0
	fa := 0

	for i := 0; i < 20; i++ {
		if c.Evaluate(model.Entity{}) {
			tr++
		} else {
			fa++
		}
	}

	assert.NotEqual(t, tr, 0)
	assert.NotEqual(t, tr, 0)
}
