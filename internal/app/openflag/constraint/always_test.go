package constraint_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/assert"
)

func TestAlwaysRule(t *testing.T) {
	rc := model.Constraint{
		Name:       constraint.AlwaysConstraintName,
		Parameters: json.RawMessage(`{}`),
	}

	err := constraint.Validate(rc.Name, rc.Parameters)
	assert.NoError(t, err)

	c, err := constraint.New(rc.Name, rc.Parameters)
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		assert.Equal(t, true, c.Evaluate(model.Entity{}))
	}
}
